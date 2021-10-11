package handler

import (
	"errors"
	"fmt"
	"sort"

	"github.com/SDJLee/mercedes-benz/model"
	"github.com/SDJLee/mercedes-benz/util"
	"gopkg.in/guregu/null.v3"
)

func computeArrival(reqBody *model.Request, transId int64) *model.Response {

	// TODO: step 1: find charge level and handle error
	// TODO: step 2: find distance and handle error
	// TODO: step 3: check if charge is required. if no, return response. else continue to step 4
	// TODO: step 4: find stations
	// TODO: step 5: find route

	// step 1: find charge level and handle error
	chargeLevel, err := getChargeLevel(reqBody)
	if err != nil || chargeLevel.Error.Valid {
		fmt.Println("error on fetching charge level", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, 0, 0, transId, true)
	}
	// if chargeLevel.Error.Valid {
	// 	fmt.Println("api error on fetching charge level", chargeLevel.Error.String)
	// 	return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, 0, 0, transId, true)
	// }
	fmt.Println("chargeLevel", chargeLevel)

	// step 2: find distance and handle error
	travelDistance, err := getTravelDistance(reqBody)
	if err != nil || travelDistance.Error.Valid {
		fmt.Println("error on fetching travel distance", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, 0, chargeLevel.CurrentChargeLevel, transId, true)
	}
	// if travelDistance.Error.Valid {
	// 	fmt.Println("api error on fetching travel distance", travelDistance.Error.String)
	// 	return generateExceptionResp(transId, true)
	// }
	fmt.Println("travelDistance", travelDistance)

	// step 3:
	// TODO: check if distance < charge left. If true, return as no charge required
	// if false, continue further.
	if chargeLevel.CurrentChargeLevel >= travelDistance.Distance {
		// with current charge level greater/equal to the total distance, there is no need to charge
		// when current charge level is equal to total distance, the charge level on arriving
		// the destination will be 0 which is acceptable.
		response := *&model.Response{
			TransactionID:      transId,
			Vin:                null.StringFrom(reqBody.Vin),
			Source:             null.StringFrom(reqBody.Source),
			Destination:        null.StringFrom(reqBody.Destination),
			CurrentChargeLevel: null.IntFrom(chargeLevel.CurrentChargeLevel),
			Distance:           null.IntFrom(travelDistance.Distance),
			IsChargingRequired: null.BoolFrom(false),
			ChargingStations:   nil,
			Errors:             nil,
		}
		fmt.Println("final response", response)
		return &response
	}

	// step 4: find stations
	chargeStations, err := getChargingStations(reqBody)
	if err != nil || chargeStations.Error.Valid {
		fmt.Println("error on fetching charging stations", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, travelDistance.Distance, chargeLevel.CurrentChargeLevel, transId, true)
	}
	// if chargeStations.Error.Valid {
	// 	fmt.Println("api error on fetching charge stations", chargeStations.Error.String)
	// 	return generateExceptionResp(transId, true)
	// }
	fmt.Println("chargeStations", chargeStations)

	stationsVisited, err := computeR(chargeStations.ChargingStations, chargeLevel.CurrentChargeLevel, travelDistance.Distance, 0)
	if err != nil {
		fmt.Println("error on fetching charging stations", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, travelDistance.Distance, chargeLevel.CurrentChargeLevel, transId, false)
	}

	sort.Strings(stationsVisited)

	// step 5: find route
	// mapper := computeDistance(chargeStations.ChargingStations, chargeLevel.CurrentChargeLevel, travelDistance.Distance, 0)
	// for primaryStation, stations := range mapper {
	// 	extdStationNames := extractStationNames(stations)
	// 	fmt.Printf("from station %v, need to visit stations %v\n", primaryStation, extdStationNames)
	// }
	// stationsToVisit := calculateMinimumDistance(chargeStations.ChargingStations, chargeLevel.CurrentChargeLevel, travelDistance.Distance, 0)

	response := *&model.Response{
		TransactionID:      transId,
		Vin:                null.StringFrom(reqBody.Vin),
		Source:             null.StringFrom(reqBody.Source),
		Destination:        null.StringFrom(reqBody.Destination),
		CurrentChargeLevel: null.IntFrom(chargeLevel.CurrentChargeLevel),
		Distance:           null.IntFrom(travelDistance.Distance),
		IsChargingRequired: null.BoolFrom(true),
		ChargingStations:   stationsVisited,
		Errors:             nil,
	}
	fmt.Println("final response", response)
	return &response
}

func getChargeLevel(reqBody *model.Request) (*model.ResChargeLevel, error) {
	chargeLevelReq := *&model.ReqChargeLevel{
		Vin: reqBody.Vin,
	}
	chargeLevel, err := GetChargeLevel(&chargeLevelReq)
	if err != nil {
		return nil, err
	}
	return chargeLevel, nil
}

func getTravelDistance(reqBody *model.Request) (*model.ResTravelDistance, error) {
	travelDistanceReq := *&model.ReqTravelDistance{
		Source:      reqBody.Source,
		Destination: reqBody.Destination,
	}
	travelDistance, err := GetTravelDistance(&travelDistanceReq)
	if err != nil {
		return nil, err
	}
	return travelDistance, nil
}

func getChargingStations(reqBody *model.Request) (*model.ResChargeStations, error) {
	chargingStationsReq := *&model.ReqChargeStations{
		Source:      reqBody.Source,
		Destination: reqBody.Destination,
	}
	chargingStations, err := GetChargingStations(&chargingStationsReq)
	if err != nil {
		return nil, err
	}
	return chargingStations, nil
}

func extractStationNames(chargingStations []*model.Station) []string {
	// TODO: perform basic null checks and other checks

	stationNames := make([]string, 0)
	for _, stations := range chargingStations {
		stationNames = append(stationNames, stations.Name)
	}
	return stationNames
}

func generateExceptionResp(vin string, source string, dest string, distance int64, chargeLevel int64, transId int64, techExp bool) *model.Response {

	generateTechException := func() []*model.ResError {
		resErrors := make([]*model.ResError, 0)
		resError := &model.ResError{
			ID:          9999,
			Description: "Technical Exception",
		}
		resErrors = append(resErrors, resError)
		return resErrors
	}

	generateUnreachableException := func() []*model.ResError {
		resErrors := make([]*model.ResError, 0)
		resError := &model.ResError{
			ID:          8888,
			Description: "unable to reach the destination with the current charge level",
		}
		resErrors = append(resErrors, resError)
		return resErrors
	}

	var errors []*model.ResError
	if techExp {
		errors = generateTechException()
	} else {
		errors = generateUnreachableException()
	}

	response := &model.Response{
		Vin:           null.StringFrom(vin),
		Source:        null.StringFrom(source),
		Destination:   null.StringFrom(dest),
		TransactionID: transId,
		Errors:        errors,
	}
	if distance > 0 {
		response.Distance = null.IntFrom(distance)
	}
	if chargeLevel > 0 {
		response.CurrentChargeLevel = null.IntFrom(chargeLevel)
	}
	return response
}

// func calculateMinimumDistance(chargingStations []*model.Station, currentCharge int, distanceToDest int, distanceTravelled int) []*model.Station {
// 	fmt.Println("\n\n---------")
// 	fmt.Printf("Calculating minimum distance with current charge %v distance to destination %v total distance travelled %v\n", currentCharge, distanceToDest, distanceTravelled)
// 	stationsToVisit := make([]*model.Station, 0)
// 	if currentCharge >= distanceToDest {
// 		fmt.Printf("EXIT CONDITION MET :: Current charge is >= to distance to destination :: %v >= %v\n", currentCharge, distanceToDest)
// 		return stationsToVisit
// 	}

// 	nextSetOfStations := extractStationNames(chargingStations)
// 	fmt.Println("visiting stations", nextSetOfStations)

// 	for stationIdx, station := range chargingStations {
// 		fmt.Printf("\tChecking stations %s \n", station.Name)
// 		// The below condition filters the station that are reachable only with current charge level.
// 		if currentCharge >= station.Distance {
// 			// Compute remaining charge. Its the difference between the current charge level and distance travelled to the station.
// 			// TODO: check if math.abs is required
// 			chargeLeft := currentCharge - (station.Distance - distanceTravelled)
// 			// Compute distance left. Its the difference between the station and destination.
// 			// TODO: check if math.abs is required
// 			distanceLeft := distanceToDest - station.Distance
// 			// The car is now refilled. Compute the new charge level after refil.
// 			refilledCharge := chargeLeft + station.Limit
// 			// // The below is the total distance travelled by the car from source.
// 			// totalDistanceTravelled := distanceTravelled + station.Distance
// 			// The below condition means, the amout refilled will be sufficient to reach the destination.
// 			fmt.Printf("\tcurrent charge - %v chargeLeft - %v distanceLeft - %v refilledCharge - %v stationDistance - %v distanceTravelled - %v \n",
// 				currentCharge, chargeLeft, distanceLeft, refilledCharge, station.Distance, distanceTravelled)
// 			if refilledCharge >= distanceLeft {
// 				// success case
// 				fmt.Printf("\t\tFinal station met :: refilledCharge >= distanceLeft :: %v >= %v \n", refilledCharge, distanceLeft)

// 				_stations := extractStationNames(stationsToVisit)
// 				fmt.Println("before appending stations", _stations)
// 				stationsToVisit = append(stationsToVisit, station)
// 				_stations = extractStationNames(stationsToVisit)
// 				fmt.Println("after appending stations", _stations)
// 				return stationsToVisit
// 			}
// 			// Check what next station is reachable
// 			fmt.Println("\t need to visit next stations")
// 			// TODO: check index out of bound
// 			return calculateMinimumDistance(chargingStations[stationIdx+1:], refilledCharge, distanceToDest, station.Distance)
// 		}
// 	}
// 	fmt.Println("---------")
// 	return stationsToVisit
// }

// func computeDistance(chargingStations []*model.Station, currentChargeLeft int, totalDistanceToDest int, totalDistanceTravelled int) map[string][]*model.Station {
// 	fmt.Println("\n\n---------")
// 	fmt.Printf("computeDistance with current charge %v distance to destination %v total distance travelled %v\n", currentChargeLeft, totalDistanceToDest, totalDistanceTravelled)
// 	if currentChargeLeft >= totalDistanceToDest {
// 		fmt.Printf("EXIT CONDITION MET :: Current charge is >= to distance to destination :: %v >= %v\n", currentChargeLeft, totalDistanceToDest)
// 		return make(map[string][]*model.Station)
// 	}
// 	// reachableStations := make([]*model.Station, 0)

// 	var computeRoutes func(stations []*model.Station, currentCharge int, distanceToDest int, distanceTravelled int) map[string][]*model.Station
// 	computeRoutes = func(stations []*model.Station, currentCharge, distanceToDest, distanceTravelled int) map[string][]*model.Station {
// 		mapper := make(map[string][]*model.Station)
// 		for _, station := range stations {
// 			fmt.Printf("\tChecking stations %s \n", station.Name)
// 			if currentCharge >= station.Distance {
// 				chargeLeft := currentCharge - (station.Distance - distanceTravelled)
// 				distanceLeft := distanceToDest - station.Distance
// 				refilledCharge := chargeLeft + station.Limit
// 				fmt.Printf("\tcurrent charge - %v chargeLeft - %v distanceLeft - %v refilledCharge - %v stationDistance - %v distanceTravelled - %v \n",
// 					currentCharge, chargeLeft, distanceLeft, refilledCharge, station.Distance, distanceTravelled)
// 				stationArr := mapper[station.Name]
// 				if stationArr == nil {
// 					stationArr = make([]*model.Station, 0)
// 				}
// 				stationArr = append(stationArr, station)
// 				mapper[station.Name] = stationArr
// 				if refilledCharge >= distanceLeft {
// 					// success
// 					fmt.Printf("\t\tFinal station met :: refilledCharge >= distanceLeft :: %v >= %v \n", refilledCharge, distanceLeft)
// 					continue
// 				}
// 				// reachableStations = append(reachableStations, station)
// 				fmt.Println("\t need to visit next stations")
// 				// return computeRoutes(stations[stationIdx+1:], refilledCharge, distanceToDest, station.Distance)
// 			} else {
// 				fmt.Printf("will not be able to reach station %v with charge %v \n", station.Name, currentCharge)
// 			}
// 		}
// 		return mapper
// 	}
// 	mapper := computeRoutes(chargingStations, currentChargeLeft, totalDistanceToDest, totalDistanceTravelled)
// 	fmt.Println("---------")
// 	return mapper
// }

func computeR(chargingStations []*model.Station, availableCharge int64, distanceToDest int64, totalDistanceTravelled int64) ([]string, error) {
	fmt.Println("\n\n---------")
	fmt.Printf("computeR with availableCharge %v distanceToDest %v totalDistanceTravelled %v\n", availableCharge, distanceToDest, totalDistanceTravelled)
	// prioritizedStations := make(map[string]int)
	// for _, station := range chargingStations {
	// 	prioritizedStations[station.Name] = station.Limit
	// }
	stationsVisited := make([]string, 0)
	pq := util.InitQueue()
	fmt.Println("first pq print")
	// pq.Print()

	for _, station := range chargingStations {
		fmt.Println("\n\n-------------")

		// fmt.Printf("outside address of pq %p\n", pq)
		fmt.Printf("checking fuel availableCharge < (station.Distance - totalDistanceTravelled) :: %v < %v - %v = %v \n", availableCharge, station.Distance, totalDistanceTravelled, (station.Distance - totalDistanceTravelled))
		for availableCharge < (station.Distance - totalDistanceTravelled) {
			fmt.Println("-------------")
			fmt.Println("\tfuel not sufficient, needs refill")
			fmt.Println("\tpq len:: ", pq.Len())
			if pq.IsEmpty() {
				// TODO: return error here
				fmt.Println("\t\tout of fuel")
				return nil, errors.New("out of fuel")
			}
			refillingStation := pq.PopItem()
			nextMaxFuel := refillingStation.Priority
			refillStationData := refillingStation.Data.(*model.Station)
			stationsVisited = append(stationsVisited, refillingStation.Value)
			chargeLeft := availableCharge - (refillStationData.Distance - totalDistanceTravelled)
			totalDistanceTravelled = refillStationData.Distance
			// fmt.Printf("\tcomputing charge left with params :: availableCharge - (refillStationData.Distance - station.Distance) = chargeLeft :: %v - (%v - %v) = %v \n", availableCharge, refillStationData.Distance, station.Distance, chargeLeft)
			fmt.Printf("\tcomputing charge left with params :: availableCharge - refillStationData.Distance = chargeLeft :: %v - %v = %v \n", availableCharge, refillStationData.Distance, chargeLeft)
			fmt.Println("\ttotalDistanceTravelled", totalDistanceTravelled)
			fmt.Printf("\trefilling at station %v \n\t\t availableCharge %v \n\t\t chargeLeft %v \n\t\t fuel at station %v \n\t\t distanceTravelled %v \n\t\t distanceToDest %v\n",
				refillingStation.Value, availableCharge, chargeLeft, refillingStation.Priority, refillStationData.Distance, distanceToDest)
			availableCharge = chargeLeft + nextMaxFuel
			fmt.Printf("\trefilled at station %v availableCharge %v with fuel %v distanceToDest %v\n", refillingStation.Value, availableCharge, refillingStation.Priority, distanceToDest)
		}
		pq.PushItem(&util.QueueItem{
			Value:    station.Name,
			Priority: station.Limit,
			Data:     station,
		})
		fmt.Printf("added station %v to queue \n", station.Name)
		// pq.Print()
		fmt.Println("visited stations: ", stationsVisited)
		fmt.Println("-------------")
	}
	fmt.Println()
	fmt.Println()

	for availableCharge < totalDistanceTravelled {
		fmt.Printf("last resort :: checking fuel availableCharge < distanceToDest :: %v < %v \n", availableCharge, distanceToDest)
		if pq.IsEmpty() {
			// TODO: return error here
			fmt.Println("last resort :: out of fuel")
			return nil, errors.New("out of fuel")
		}
		fmt.Printf("Available stations: ")
		// pq.Print()
		refillingStation := pq.PopItem()
		nextMaxFuel := refillingStation.Priority
		// refillStationData := refillingStation.Data.(*model.Station)
		stationsVisited = append(stationsVisited, refillingStation.Value)
		fmt.Printf("refilling last resort at station %v availableCharge %v with fuel %v distanceToDest %v\n", refillingStation.Value, availableCharge, refillingStation.Priority, distanceToDest)
		availableCharge = availableCharge + nextMaxFuel
		fmt.Printf("refilled last resort at station %v availableCharge %v with fuel %v distanceToDest %v\n", refillingStation.Value, availableCharge, refillingStation.Priority, distanceToDest)
	}
	// for pq.Len() > 0 {
	// 	item := pq.PopItem()
	// 	fmt.Printf("%.2d:%s \n", item.Priority, item.Value)
	// }
	// pq.Print()
	fmt.Println("---------")
	fmt.Println("stationsVisited", stationsVisited)
	return stationsVisited, nil
}
