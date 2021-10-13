package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SDJLee/mercedes-benz/model"
	"github.com/SDJLee/mercedes-benz/util"
	"github.com/spf13/viper"
)

const (
	ReqPost      = "POST"
	reqTestCase1 = "{ \"vin\": \"W1K2062161F0046\", \"source\": \"Home\", \"destination\": \"Movie Theatre\" }"
	reqTestCase2 = "{ \"vin\": \"W1K2062161F0080\", \"source\": \"Home\", \"destination\": \"Airport\" }"
	reqTestCase3 = "{ \"vin\": \"W1K2062161F0080\", \"source\": \"@$%%%\", \"destination\": \"Airport\" }"
	reqTestCase4 = "{ \"vin\": \"W1K2062161F0046\", \"source\": \"Home\", \"destination\": \"Movie Theatre\" }"
)

// To test health endpoint
func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := executeRequest(req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"breathing..."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//{ "transactionId": 2, "vin": "W1K2062161F0033", "source": "Home", "destination": "Lake", "distance": 20, "currentChargeLevel": 93, "isChargingRequired": false }
func TestCase1(t *testing.T) {

	responseBody, err := performApiCall(reqTestCase1, t)
	if err != nil {
		t.Fatal(err)
	}
	reqBody := &model.Request{}
	err = json.Unmarshal([]byte(reqTestCase1), reqBody)
	if err != nil {
		t.Fatal(err)
	}

	if responseBody.Vin.String != reqBody.Vin {
		t.Errorf("invalid 'vin' in response. It should be '%v' but it is '%v'", reqBody.Vin, responseBody.Vin.String)
	}

	if responseBody.Source.String != reqBody.Source {
		t.Errorf("invalid 'source' in response. It should be '%v' but it is '%v'", reqBody.Source, responseBody.Source.String)
	}

	if responseBody.Destination.String != reqBody.Destination {
		t.Errorf("invalid 'destination' in response. It should be '%v' but it is '%v'", reqBody.Destination, responseBody.Destination.String)
	}

	if responseBody.IsChargingRequired.Bool {
		t.Errorf("charging required should be false")
	}
}

//{ "transactionId": 3, "vin": "W1K2062161F0080", "source": "Home", "destination": "Airport", "distance": 60, "currentChargeLevel": 2, "isChargingRequired": true, "errors": [ { "id": 8888, "description": "Unable to reach the destination with the current fuel level" } ] }
func TestCase2(t *testing.T) {

	responseBody, err := performApiCall(reqTestCase2, t)
	if err != nil {
		t.Fatal(err)
	}
	reqBody := &model.Request{}
	err = json.Unmarshal([]byte(reqTestCase2), reqBody)
	if err != nil {
		t.Fatal(err)
	}

	if responseBody.Vin.String != reqBody.Vin {
		t.Errorf("invalid 'vin' in response. It should be '%v' but it is '%v'", reqBody.Vin, responseBody.Vin.String)
		t.Fail()
	}

	if responseBody.Source.String != reqBody.Source {
		t.Errorf("invalid 'source' in response. It should be '%v' but it is '%v'", reqBody.Source, responseBody.Source.String)
		t.Fail()
	}

	if responseBody.Destination.String != reqBody.Destination {
		t.Errorf("invalid 'destination' in response. It should be '%v' but it is '%v'", reqBody.Destination, responseBody.Destination.String)
		t.Fail()
	}

	if responseBody.Errors == nil {
		t.Error("errors should not be nil")
		t.FailNow()
	}

	if len(responseBody.Errors) <= 0 {
		t.Errorf("error with ID %v should have returned", util.ErrUnreachableId)
		t.FailNow()
	}

	if len(responseBody.Errors) > 1 {
		t.Fatal("there shouldn't be more than one error")
		t.Fail()
	}

	if responseBody.Errors[0].ID != util.ErrUnreachableId {
		t.Errorf("expected error with ID %v but got %v", util.ErrUnreachableId, responseBody.Errors[0].ID)
		t.Fail()
	}

	if responseBody.Errors[0].Description != util.ErrUnreachableMsg {
		t.Error("invalid error description")
		t.Fail()
	}

}

// {"transactionId": 1,"errors": [{"id": 9999,"description": "Technical Exception"}]}
func TestCase3(t *testing.T) {

	responseBody, err := performApiCall(reqTestCase3, t)
	if err != nil {
		t.Fatal(err)
	}
	reqBody := &model.Request{}
	err = json.Unmarshal([]byte(reqTestCase3), reqBody)
	if err != nil {
		t.Fatal(err)
	}

	if responseBody.Vin.String != reqBody.Vin {
		t.Errorf("invalid 'vin' in response. It should be '%v' but it is '%v'", reqBody.Vin, responseBody.Vin.String)
		t.Fail()
	}

	if responseBody.Source.String != reqBody.Source {
		t.Errorf("invalid 'source' in response. It should be '%v' but it is '%v'", reqBody.Source, responseBody.Source.String)
		t.Fail()
	}

	if responseBody.Destination.String != reqBody.Destination {
		t.Errorf("invalid 'destination' in response. It should be '%v' but it is '%v'", reqBody.Destination, responseBody.Destination.String)
		t.Fail()
	}

	if responseBody.Errors == nil {
		t.Error("errors should not be nil")
		t.FailNow()
	}

	if len(responseBody.Errors) <= 0 {
		t.Errorf("error with ID %v should have returned", util.ErrTechExpId)
		t.FailNow()
	}

	if len(responseBody.Errors) > 1 {
		t.Fatal("there shouldn't be more than one error")
		t.Fail()
	}

	if responseBody.Errors[0].ID != util.ErrTechExpId {
		t.Errorf("expected error with ID %v but got %v", util.ErrTechExpId, responseBody.Errors[0].ID)
		t.Fail()
	}

	if responseBody.Errors[0].Description != util.ErrTechExpMsg {
		t.Error("invalid error description")
		t.Fail()
	}

}

//{ "transactionId": 5, "vin": "W1K2062161F0046", "source": "Home", "destination": "Movie Theatre", "distance": 50, "currentChargeLevel": 17, "isChargingRequired": true, "chargingStations": [ { "name": "S1", "distance": 10, "limit": 20 }, { "name": "S2", "distance": 25, "limit": 15 } ] }
func TestCase4(t *testing.T) {

	responseBody, err := performApiCall(reqTestCase4, t)
	if err != nil {
		t.Fatal(err)
	}
	reqBody := &model.Request{}
	err = json.Unmarshal([]byte(reqTestCase4), reqBody)
	if err != nil {
		t.Fatal(err)
	}

	if responseBody.Vin.String != reqBody.Vin {
		t.Errorf("invalid 'vin' in response. It should be '%v' but it is '%v'", reqBody.Vin, responseBody.Vin.String)
		t.Fail()
	}

	if responseBody.Source.String != reqBody.Source {
		t.Errorf("invalid 'source' in response. It should be '%v' but it is '%v'", reqBody.Source, responseBody.Source.String)
		t.Fail()
	}

	if responseBody.Destination.String != reqBody.Destination {
		t.Errorf("invalid 'destination' in response. It should be '%v' but it is '%v'", reqBody.Destination, responseBody.Destination.String)
		t.Fail()
	}

	if responseBody.Errors != nil {
		t.Error("errors should be nil")
		t.Fail()
	}

	if !responseBody.IsChargingRequired.Bool {
		t.Error("charge required should be true")
		t.Fail()
	}

	if responseBody.ChargingStations == nil {
		t.Fatal("charging station shouldn't be nil")
		t.FailNow()
	}

	if len(responseBody.ChargingStations) <= 0 {
		t.Fatal("charging station shouldn't be empty")
		t.FailNow()
	}

	if len(responseBody.ChargingStations) < 2 {
		t.Error("this testcase should return 2 charging stations")
		t.Fail()
	}

	if responseBody.ChargingStations[0] != "S1" || responseBody.ChargingStations[1] != "S2" {
		t.Error("this testcase should return S1 and S2 for charging stations")
		t.Fail()
	}
}

func TestCase5(t *testing.T) {
	everyStation := make([]*model.Station, 4)
	everyStation[0] = &model.Station{
		Name:     "S1",
		Limit:    20,
		Distance: 10,
	}
	everyStation[1] = &model.Station{
		Name:     "S2",
		Limit:    15,
		Distance: 25,
	}
	everyStation[2] = &model.Station{
		Name:     "S3",
		Limit:    10,
		Distance: 33,
	}
	everyStation[3] = &model.Station{
		Name:     "S4",
		Limit:    10,
		Distance: 40,
	}

	stationsVisited, err := computeRoute(everyStation, 17, 50)

	if err != nil {
		t.Error("test case shouldn't return error")
		t.FailNow()
	}

	if stationsVisited == nil {
		t.Fatal("charging station shouldn't be nil")
		t.FailNow()
	}

	if len(stationsVisited) <= 0 {
		t.Fatal("charging station shouldn't be empty")
		t.FailNow()
	}

	if len(stationsVisited) != 2 {
		t.Error("this testcase should return 2 charging stations")
		t.Fail()
	}

	if stationsVisited[0] != "S1" || stationsVisited[1] != "S2" {
		t.Error("this testcase should return S1 and S2 for charging stations")
		t.Fail()
	}
}

func TestCase6(t *testing.T) {
	everyStation := make([]*model.Station, 4)
	everyStation[0] = &model.Station{
		Name:     "S1",
		Limit:    20,
		Distance: 10,
	}
	everyStation[1] = &model.Station{
		Name:     "S2",
		Limit:    30,
		Distance: 25,
	}
	everyStation[2] = &model.Station{
		Name:     "S3",
		Limit:    45,
		Distance: 33,
	}
	everyStation[3] = &model.Station{
		Name:     "S4",
		Limit:    20,
		Distance: 40,
	}

	stationsVisited, err := computeRoute(everyStation, 17, 90)

	if err != nil {
		t.Error("test case shouldn't return error")
		t.FailNow()
	}

	if stationsVisited == nil {
		t.Fatal("charging station shouldn't be nil")
		t.FailNow()
	}

	if len(stationsVisited) <= 0 {
		t.Fatal("charging station shouldn't be empty")
		t.FailNow()
	}

	if len(stationsVisited) != 4 {
		t.Error("this testcase should return 4 charging stations")
		t.Fail()
	}

	if stationsVisited[0] != "S1" || stationsVisited[1] != "S3" {
		t.Error("this testcase should return S1 and S3 for charging stations")
		t.Fail()
	}
}

func TestCaseInvalidReq(t *testing.T) {
	req, err := http.NewRequest(ReqPost, computeBaseUrl(util.ApiComputeRoute), strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	rr := executeRequest(req)
	t.Log("code:: ", rr.Code)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

// helper method to make HTTP request
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func computeBaseUrl(apiName string) string {
	return fmt.Sprintf("http://localhost:%s%s%s%s", viper.GetString(util.Port), util.ApiBasePath, util.ApiV1, apiName)
}

func performApiCall(requestPayload string, t *testing.T) (*model.Response, error) {
	bufferPayload := bytes.NewBuffer([]byte(requestPayload))
	url := computeBaseUrl(util.ApiComputeRoute)
	req, err := http.NewRequest(ReqPost, url, bufferPayload)
	if err != nil {
		return nil, err
	}

	for key, val := range defaultHeaders {
		req.Header.Add(key, val)
	}

	rr := executeRequest(req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("invalid response: got %v want %v",
			status, http.StatusOK)
	}

	responseByte, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		return nil, err
	}
	responseBody := &model.Response{}
	err = json.Unmarshal(responseByte, responseBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
