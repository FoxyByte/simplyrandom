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

func (r *random) request(params map[string]interface{}) (map[string]interface{}, error) {
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
