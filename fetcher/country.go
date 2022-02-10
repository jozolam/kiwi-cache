package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Locations []struct {
		IntID int    `json:"int_id"`
		ID    string `json:"id"`
	} `json:"locations"`
}

type CountryFetcher struct {
}

func (*CountryFetcher) FetchAll() (map[int]string, error) {
	values := make(map[int]string)
	response, err := http.Get("https://api.skypicker.com/locations/graphql?query=%7Bdump%20%28options%3A%20%7Blocation_types%3A%20%5B%22airport%22%5D%2C%20active_only%3A%20%22true%22%7D%29%20%7Bint_id%20id%7D%7D")
	if err != nil {
		return values, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return values, err
	}

	var result Response
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return values, err
	}

	for _, v := range result.Locations {
		values[v.IntID] = v.ID
	}

	return values, nil
}

func (*CountryFetcher) Fetch(id int) (string, error) {
	result, err := getResponse(fmt.Sprintf("https://api.skypicker.com/locations?type=general&key=int_id&value=%v", id))
	if err != nil {
		return "", err
	}

	if len(result.Locations) == 0 {
		return "", errors.New("location does not exists")
	}

	return result.Locations[0].ID, nil
}

func getResponse(url string) (Response, error) {
	var result Response
	response, err := http.Get(url)
	if err != nil {
		return result, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
