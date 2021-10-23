package handler

import (
	"bytes"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SDJLee/mercedes-benz/metrics"
	"github.com/SDJLee/mercedes-benz/model"
	"github.com/SDJLee/mercedes-benz/util"
	"github.com/spf13/viper"
)

var defaultHeaders = map[string]string{
	"Content-Type":  "application/json",
	"Response-Type": "application/json",
}

// retrieves current charge level
func GetChargeLevel(requestBody *model.ReqChargeLevel) (*model.ResChargeLevel, error) {
	defer metrics.StatTime("api.chargelevel")()
	logger.Info("retrieving charge level data")
	defer logger.Info("retrieved charge level data")
	url := fmt.Sprintf("%s/charge_level", viper.GetString(util.ApiAddress))
	logger.Debugf("API url to retrieve charge level: %s", url)

	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	responseByte, err := makePostRequest(url, jsonPayload)
	if err != nil {
		return nil, err
	}

	response := &model.ResChargeLevel{}
	err = json.Unmarshal(responseByte, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// retrieves travel distance
func GetTravelDistance(requestBody *model.ReqTravelDistance) (*model.ResTravelDistance, error) {
	defer metrics.StatTime("api.traveldistance")()
	logger.Info("retrieving travel distance data")
	defer logger.Info("retrieved travel distance data")
	url := fmt.Sprintf("%s/distance", viper.GetString(util.ApiAddress))
	logger.Debugf("API url to retrieve travel distance: %s", url)

	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	responseByte, err := makePostRequest(url, jsonPayload)
	if err != nil {
		return nil, err
	}

	response := &model.ResTravelDistance{}
	err = json.Unmarshal(responseByte, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// retrieves charging stations
func GetChargingStations(requestBody *model.ReqChargeStations) (*model.ResChargeStations, error) {
	defer metrics.StatTime("api.chargestation")()
	logger.Info("retrieving charge stations data")
	defer logger.Info("retrieved charge stations data")
	url := fmt.Sprintf("%s/charging_stations", viper.GetString(util.ApiAddress))
	logger.Debugf("API url to retrieve charge stations: %s", url)

	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	responseByte, err := makePostRequest(url, jsonPayload)
	if err != nil {
		return nil, err
	}

	response := &model.ResChargeStations{}
	err = json.Unmarshal(responseByte, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// common method to perform http post request
func makePostRequest(url string, bytePayload []byte) ([]byte, error) {
	bufferPayload := bytes.NewBuffer(bytePayload)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bufferPayload)
	if err != nil {
		return nil, err
	}

	for key, val := range defaultHeaders {
		request.Header.Add(key, val)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseByte, nil
}
