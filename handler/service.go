package handler

import (
	"errors"
	"sort"

	"github.com/SDJLee/mercedes-benz/metrics"
	"github.com/SDJLee/mercedes-benz/model"
	"github.com/SDJLee/mercedes-benz/util"
	"gopkg.in/guregu/null.v3"
)

// TODO: logger to file with rotation
// TODO: dockerize
// TODO: test case
// TODO: PPT
// TODO: Documentation comments
// TODO: clean up
//		-- replace fmt with logger
//		-- remove unwanted fmt
//		-- remove unwanted comments
// TODO: handle panic
// TODO: recheck computation logic. PQ changes order of stations. will the distance be valid if order of stations are changed?

// computeArrival contains the logic that handles http calls to retrieve charge level, distance to destination and charging station data to determine
// if the car can travel to destination with current charge level. If the car cannot reach the destination with current charge level,
// the logic computes the minimum number of charging stations to visit.
// It returns the response that contains the cummulative information from above API calls and computed stations to visit list. In case of error or if
// the destination/station cannot be reached with current charge, it returns appropriate error code and message.
func computeArrival(reqBody *model.Request, transId int64) (response *model.Response) {
	// recover a panic and return technical exception
	defer func() {
		if ex := recover(); ex != nil {
			logger.Error("panic recovered", ex)
			response = generateExceptionResp("", "", "", 0, 0, transId, true)
		}
	}()

	// step 1: find charge level and handle error
	chargeLevel, err := getChargeLevel(reqBody)
	if err != nil || chargeLevel.Error.Valid {
		logger.Error("error on fetching charge level", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, 0, 0, transId, true)
	}
	logger.Debug("chargeLevel", chargeLevel)

	// step 2: find distance and handle error
	travelDistance, err := getTravelDistance(reqBody)
	if err != nil || travelDistance.Error.Valid {
		logger.Error("error on fetching travel distance", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, 0, chargeLevel.CurrentChargeLevel, transId, true)
	}
	logger.Debug("travelDistance", travelDistance)

	// step 3: handle if current level is sufficient to reach the destination
	if chargeLevel.CurrentChargeLevel >= travelDistance.Distance {
		// with current charge level greater/equal to the total distance, there is no need to charge
		// when current charge level is equal to total distance, the charge level on arriving
		// the destination will be 0 which is acceptable.
		response = &model.Response{
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
		logger.Debug("final response", response)
		return response
	}

	// at this point, we know that with current charge level, we cannot reach the distance. continue further to retrieve list of available charging
	// stations and pick the minimum number of stations to visit.

	// step 4: find stations
	chargeStations, err := getChargingStations(reqBody)
	if err != nil || chargeStations.Error.Valid {
		logger.Error("error on fetching charging stations", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, travelDistance.Distance, chargeLevel.CurrentChargeLevel, transId, true)
	}
	logger.Debug("chargeStations", chargeStations)

	// step 5: compute the minimum number of stations to visit.
	stationsVisited, err := computeRoute(chargeStations.ChargingStations, chargeLevel.CurrentChargeLevel, travelDistance.Distance)
	if err != nil {
		logger.Error("error on fetching charging stations", err)
		return generateExceptionResp(reqBody.Vin, reqBody.Source, reqBody.Destination, travelDistance.Distance, chargeLevel.CurrentChargeLevel, transId, false)
	}

	// sort the stations slice order the station names lexicographically
	sort.Strings(stationsVisited)

	response = &model.Response{
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
	logger.Debug("final response", response)
	return response
}

// getChargeLevel method handles the API call to retrieve current charge level
func getChargeLevel(reqBody *model.Request) (*model.ResChargeLevel, error) {
	defer metrics.MonitorTimeElapsed("constructing charge level")()
	chargeLevelReq := *&model.ReqChargeLevel{
		Vin: reqBody.Vin,
	}
	chargeLevel, err := GetChargeLevel(&chargeLevelReq)
	if err != nil {
		return nil, err
	}
	return chargeLevel, nil
}

// getTravelDistance method handles the API call to retrieve the travel distance
func getTravelDistance(reqBody *model.Request) (*model.ResTravelDistance, error) {
	defer metrics.MonitorTimeElapsed("constructing travel distance")()
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

// getChargingStations method handles the API call to retrieve slice of charging stations between source and destination
func getChargingStations(reqBody *model.Request) (*model.ResChargeStations, error) {
	defer metrics.MonitorTimeElapsed("constructing charge stations")()
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

// generateExceptionResp is a helper method to generate error responses. The type of error is differentiated by techExp param.
// If techExp is true, error 9999 is generated. Else error 8888 is generated.
func generateExceptionResp(vin string, source string, dest string, distance int64, chargeLevel int64, transId int64, techExp bool) *model.Response {

	// generates an error for "Technical Exception" for invalid request/data or computational failure
	generateTechException := func() []*model.ResError {
		resErrors := make([]*model.ResError, 0)
		resError := &model.ResError{
			ID:          9999,
			Description: "Technical Exception",
		}
		resErrors = append(resErrors, resError)
		return resErrors
	}

	// generates an error to denote that the car will not be able to reach the destination due to insufficient charge
	generateUnreachableException := func() []*model.ResError {
		resErrors := make([]*model.ResError, 0)
		resError := &model.ResError{
			ID:          8888,
			Description: "Unable to reach the destination with the current charge level",
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

	// these are added since the sample responses in the problem statement has them. If the values are 0, set null in response.
	if distance > 0 {
		response.Distance = null.IntFrom(distance)
	}
	if chargeLevel > 0 {
		response.CurrentChargeLevel = null.IntFrom(chargeLevel)
	}
	return response
}

// computeRoute computes the slice of minimum number of stations to be visited to recharge before reaching the destinatio.
// The logic follows a greedy approach where we charge the car only at stations that can provide maximum number of charges when compared to all other stations at that state.
// Below is the logical explanation of the method.
// 1. We find out the maximum distance the car can travel with available charge.
// 2. If the destination can be reached with available charge, we return an empty slice. It means that there is no need for the car to stop for recharging since the available charge is sufficient.
// 3. If the destination cannot be reached with available charge, the logic simulates the car to travel to maximum distance possible noting down the stations along the route in priority queue.
// The priority queue will be in descending order respect to the charge available in station. For example, if the station and charge pair are S1:10, S2:20, S3:30, then the priority queue will return
// in the order S3:30, S2:20, S1:10. We always pick the next station that provides maximum charge.
// 4. If the charge in a station is not sufficient, we pick the next station from the priority queue. This is done till either the queue is empty or the charge becomes sufficient.
// 5. The station names where the car recharges are added to the returning slice.
// The method returns a slice containing names of the stations where the car is recharged. The slice is empty if no station is visited. This is when the charge is sufficient to reach destination.
// The method also returns a error variable. This error is to denote that the car will not make it to the destination as there is no sufficient charge.
// The time complexity of this logic is O(nlog(n)). We iterate n times and greedily check if recharge is required.
// The space complexity of this logic is O(n)
func computeRoute(chargingStations []*model.Station, availableCharge int64, distanceToDest int64) ([]string, error) {
	logger.Info("computing route")
	defer logger.Info("route computed")
	var distanceTravelled int64 = 0
	logger.Debug("\n\n---------")
	logger.Debugf("computeR with availableCharge %v distanceToDest %v distanceTravelled %v\n", availableCharge, distanceToDest, distanceTravelled)
	stationsVisited := make([]string, 0)
	pq := util.InitQueue()

	// if available charge is >= distance to destination, there is no need to stop at stations to recharge. Return empty slice.
	if availableCharge >= distanceToDest {
		return stationsVisited, nil
	}

	// iterate through stations to apply greedy approach
	for _, station := range chargingStations {
		logger.Debug("\n\n-------------")

		logger.Debugf("checking charge availableCharge < (station.Distance - distanceTravelled) :: %v < %v - %v = %v \n", availableCharge, station.Distance, distanceTravelled, (station.Distance - distanceTravelled))
		// This condition is to check if the charge left in car is sufficient to reach the next station.
		// The calculation (station.Distance - distanceTravelled) is done because, the distance provided in station struct doesn't denote the distance between the stations. It denotes the distance between
		// the source and the station. As the car travels from source to destination and passes through each station, we should subtract this travelled distance with the distance provided in station struct.
		// In simpler terms, the distance provided in station struct includes the distance travelled by the car. We need the difference between the two to check if the car will be able to reach the next station
		// from a previous station/source with current charge.
		for availableCharge < (station.Distance - distanceTravelled) {
			logger.Debug("-------------")
			logger.Debug("\tcharge not sufficient, needs refill")
			logger.Debug("\tpq len:: ", pq.Len())
			// If there are no more stations left with charge, then there is no sufficient charge for the car to reach the destination. return error.
			if pq.IsEmpty() {
				logger.Warn("\t\tout of charge")
				return nil, errors.New("out of charge")
			}
			// The priority queue pops the element with greater priority value. Here, the charge left in a station is the priority. This pop will return next station that has the maximum charge left for consumption.
			refillingStation := pq.PopItem()
			nextMaxCharge := refillingStation.Priority
			refillStationData := refillingStation.Data.(*model.Station)
			// push the station into stationsVisited slice. This slice keeps track of stations that are visited to recharge.
			stationsVisited = append(stationsVisited, refillingStation.Value)
			// compute chargeLeft. It is the difference between the initial available charge from source or last station visit and the distance travelled to this station from source or a previous station.
			// (refillStationData.Distance - distanceTravelled) gives the distance between the station.
			chargeLeft := availableCharge - (refillStationData.Distance - distanceTravelled)
			// update the total distance travelled with the station's distance. This is because, the station's distance is the distance from the source.
			// Since the stations are out of order in the priority queue, the distanceTravelled should be updated only if it is lesser than the distance to the station.
			// For example, distance from source could be more for S3 when compared to S2 in the data S3:50, S2:40. But if S2 provides more charge than S3, then the priority queue
			// will pop S2 first. Without this condition, the distanceTravelled can also decrease which we should avoid. DistanceTravelled should always contain the farthest distance covered by the car.
			if distanceTravelled < refillStationData.Distance {
				distanceTravelled = refillStationData.Distance
			}
			logger.Debugf("\tcomputing charge left with params :: availableCharge - refillStationData.Distance = chargeLeft :: %v - %v = %v \n", availableCharge, refillStationData.Distance, chargeLeft)
			logger.Debug("\tdistanceTravelled", distanceTravelled)
			logger.Debugf("\trefilling at station %v \n\t\t availableCharge %v \n\t\t chargeLeft %v \n\t\t charge at station %v \n\t\t distanceTravelled %v \n\t\t distanceToDest %v\n",
				refillingStation.Value, availableCharge, chargeLeft, refillingStation.Priority, refillStationData.Distance, distanceToDest)
			// This shows the refilling process. The availableCharge value is updated to the sum between chargeLeft and the charge available at the station depicted by variable nextMaxCharge.
			availableCharge = chargeLeft + nextMaxCharge
			logger.Debugf("\trefilled at station %v availableCharge %v with charge %v distanceToDest %v\n", refillingStation.Value, availableCharge, refillingStation.Priority, distanceToDest)
		}
		// regardless if car stops for recharge, push the station into priority queue on each iteration. This station will be consumed in above for loop when charge is required.
		pq.PushItem(&util.QueueItem{
			Value:    station.Name,
			Priority: station.Limit,
			Data:     station,
		})
		logger.Debugf("added station %v to queue \n", station.Name)
		logger.Debug("visited stations: ", stationsVisited)
		logger.Debug("-------------")
	}

	// handling edge case where the car hasn't reached the destination but still has stations left to recharge.
	// The logic repeats the same step from above.
	for availableCharge < distanceTravelled {
		// TODO: Recheck this computation, move to closure
		logger.Debugf("last resort :: checking charge availableCharge < distanceToDest :: %v < %v \n", availableCharge, distanceToDest)
		if pq.IsEmpty() {
			// If there are no more stations left with charge, then there is no sufficient charge for the car to reach the destination. return error.
			logger.Debug("last resort :: out of charge")
			return nil, errors.New("out of charge")
		}
		logger.Debugf("Available stations: ")
		refillingStation := pq.PopItem()
		nextMaxCharge := refillingStation.Priority
		refillStationData := refillingStation.Data.(*model.Station)
		stationsVisited = append(stationsVisited, refillingStation.Value)
		chargeLeft := availableCharge - (refillStationData.Distance - distanceTravelled)
		distanceTravelled = refillStationData.Distance
		logger.Debugf("refilling last resort at station %v availableCharge %v with charge %v distanceToDest %v\n", refillingStation.Value, availableCharge, refillingStation.Priority, distanceToDest)
		availableCharge = chargeLeft + nextMaxCharge
		logger.Debugf("refilled last resort at station %v availableCharge %v with charge %v distanceToDest %v\n", refillingStation.Value, availableCharge, refillingStation.Priority, distanceToDest)
	}
	logger.Debug("---------")
	logger.Debug("stationsVisited", stationsVisited)
	return stationsVisited, nil
}
