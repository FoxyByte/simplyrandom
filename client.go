package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pborman/uuid"
)

const (
	baseURL = "https://api.random.org/json-rpc/2/invoke"
	errAPI  = "API Error Code %v: %q."
)

var (
	errAPIKey = errors.New("missing API Key, please provide key")
	errParams = errors.New("invalid parameter range l or r is bigger then 10")
	errJSON   = errors.New("invalid JSON")
)

//Random for new Client.
type Random struct {
	apiKey string
	client *http.Client
}

// NewRandom creates a client with the given apiKey.
func NewRandom(apiKey string) *Random {
	if apiKey == "" {
		panic(errAPIKey)
	}

	random := Random{
		apiKey: apiKey,
		client: &http.Client{},
	}

	return &random
}

func (r *Random) jsonMap(json map[string]interface{}, key string) (map[string]interface{}, error) {
	value := json[key]
	if value == nil {
		return nil, errJSON
	}

	newMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, errJSON
	}

	return newMap, nil
}

func (r *Random) request(params map[string]interface{}) (map[string]interface{}, error) {

	params["apiKey"] = r.apiKey
	getUUID := uuid.NewUUID().String()

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "generateIntegers",
		"params":  params,
		"id":      getUUID,
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestBodyReader := bytes.NewReader(requestBodyJSON)

	req, err := http.NewRequest("POST", baseURL, requestBodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	responseBody := make(map[string]interface{})
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		if len(body) > 0 {
			err = errors.New(string(body))
		}

		return nil, err
	}

	result, err := r.jsonMap(responseBody, "result")
	if err != nil {
		error, err := r.jsonMap(responseBody, "error")
		if err != nil {
			return nil, err
		}

		errorCode, _ := error["code"]
		errorMessage, _ := error["message"]
		err = fmt.Errorf(errAPI, errorCode, errorMessage)
		return nil, err
	}

	return result, nil
}

func (r *Random) req(params map[string]interface{}) ([]interface{}, error) {
	result, err := r.request(params)
	if err != nil {
		return nil, err
	}

	random, err := r.jsonMap(result, "random")
	if err != nil {
		return nil, err
	}

	data := random["data"].([]interface{})

	return data, nil
}

func (r *Random) getIntegers(l int) ([]int64, error) {
	if l < 1 || l > 10 {
		return nil, errParams
	}

	params := map[string]interface{}{
		"n":   l,
		"min": 0,
		"max": 10,
	}

	values, err := r.req(params)
	if err != nil {
		return nil, err
	}

	ints := make([]int64, len(values))
	for i, value := range values {
		f := value.(float64)
		ints[i] = int64(f)
	}

	return ints, nil
}
