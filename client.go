package main

import (
	"errors"
	"net/http"
)

const (
	baseURL = "https://api.random.org/json-rpc/2/invoke"
	errAPI  = "API Error Code %v: %q."
)

var (
	errAPIKey = errors.New("missing API Key, please provide key")
	errParams = errors.New("invalid parameter range l is bigger then 10")
	errJSON   = errors.New("invalid JSON")
)

type random struct {
	apiKey string
	client *http.Client
}

func newRandom(apiKey string) *random {
	if apiKey == "" {
		panic(errAPIKey)
	}

	random := random{
		apiKey: apiKey,
		client: &http.Client{},
	}

	return &random
}

func (r *random) req(params map[string]interface{}) ([]interface{}, error) {
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

func (r *random) getIntegers(l int) ([]int64, error) {
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
