package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var apiKey = os.Getenv("APIKEY")

func getRandoms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	length := query.Get("length")
	requests := query.Get("requests")

	len, err := strconv.Atoi(length)
	if err != nil {
		fmt.Println(err)
	}
	req, err := strconv.Atoi(requests)
	if err != nil {
		fmt.Println(err)
	}

	output := `
	{
	"stddev": %f
	"data": %d 
	}`

	var valuesAll []int64
	wg := sync.WaitGroup{}

	if req <= 10 {
		for i := 0; i < req; i++ {
			wg.Add(1)
			go func(len int) {
				random := newRandom(apiKey)
				value, err := random.getIntegers(len)
				if err != nil {
					fmt.Println(err)
				}
				stdDev := standardDev(value)
				valuesAll = append(valuesAll, value...)
				w.Write([]byte(fmt.Sprintf(output, stdDev, value)))
				wg.Done()
			}(len)
		}
		wg.Wait()
		stdDevAll := standardDev(valuesAll)
		w.Write([]byte(fmt.Sprintf(output, stdDevAll, valuesAll)))
	} else {
		w.Write([]byte(fmt.Sprintf("Sorry but 10 requests in the max")))
	}
}
