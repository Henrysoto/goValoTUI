package main

import (
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"strconv"
)

const (
	BaseAPI                        = "https://api.henrikdev.xyz/valorant/"
	EndpointLifetimeMatchesByPuuid = "v1/by-puuid/lifetime/matches/eu/%s?mode=competitive"
	EndpointLifeTimeMatchesByName  = "v1/lifetime/matches/eu/%s/%s?mode=competitive"
	EndpointMatchId                = "v2/match/%s"
	EndpointMatchesByName          = "v3/matches/eu/%s/%s"
	EndpointAccount                = "v1/account/%s/%s"
	EndpointRank                   = "v1/mmr-history/eu/%s/%s"
)

var RequestErrors = map[string]string{
	"400": "client request error",
	"403": "connection forbidden to Riot API",
	"404": "entity not found (private profile or incorrect username)",
	"408": "request timeout",
	"429": "requests rate limit exceeded (retry later)",
	"503": "Riot API is down (retry later)",
}

// Retrieve JSON data from API Endpoints
func GetData[T Account | AccountRank | MatchesData](endpoint string) (T, error) {
	var client http.Client
	var result T

	// Creating request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	req.Header.Add("accept", "application/json")

	// Retrieve response from request
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// check if status is healthy (aka 200)
	if RequestErrors[strconv.Itoa(resp.StatusCode)] != "0" && resp.StatusCode != 200 {
		return result, fmt.Errorf(RequestErrors[strconv.Itoa(resp.StatusCode)])
	}

	// Unmarshal JSON body to struct
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Return Image object from URL
func getImageFromURL(URL string) (image.Image, error) {
	var img image.Image
	data, err := http.Get(URL)
	if err != nil {
		return img, err
	}
	defer data.Body.Close()

	img, _, err = image.Decode(data.Body)
	if err != nil {
		return img, err
	}

	return img, nil
}
