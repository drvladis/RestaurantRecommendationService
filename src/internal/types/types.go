package types

import "encoding/json"

type Place struct {
	ID       int    `json:"id" csv:"id"`
	Name     string `json:"name" csv:"name"`
	Address  string `json:"address" csv:"address"`
	Phone    string `json:"phone" csv:"phone"`
	Location struct {
		Longitude float64 `json:"lon" csv:"lon"`
		Latitude  float64 `json:"lat" csv:"lat"`
	} `json:"location"`
}

type EsResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Store interface {
	// Returns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit, offset int, index string) ([]Place, int, error)
	// Returns the list of the closest places of set lat and lon in amount of (var: closest) and (or) an error in case of one
	GetClosestPlaces(lat, lon float64, index string) ([]Place, error)
}
