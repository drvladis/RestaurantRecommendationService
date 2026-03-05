package dataloader

import (
	"APIforElasticBD/internal/types"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	logger "github.com/drvladis/logging/log"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/elastic/go-elasticsearch/v9/esutil"
)

const (
	closest int = 3
	timeout int = 5
)

// Struct with GetPlaces method for separate logic and interact
type Client struct {
	Client *elasticsearch.Client
}

// Create a new elasticsearch client and returns pointer to Client struct and error
func Init() (*Client, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Creating client error:", err)
	}
	return &Client{client}, err
}

// Loads data from filename to index through esutil.NewBulkIndexer
func DataLoader(filename, index string) {
	client, err := Init()

	req, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  index,
		Client: client.Client,
	})
	if err != nil {
		log.Fatal("dataloader: Creating indexer error", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("dataloader: opening csv-file error:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	_, err = reader.Read()
	if err != nil {
		log.Fatal("dataloader: reading csv-file error:", err)
	}

	for {
		record, err := reader.Read()
		if err != nil && err != io.EOF {
			logger.Logging(fmt.Errorf("dataloader: csv.reader: %s\n", err.Error()))
			log.Fatal("dataloader: csv.reader: ", err)
		}
		if err == io.EOF {
			fmt.Printf("The document %s has been read.\n", file.Name())
			break
		}

		place := &types.Place{}
		place.ID, err = strconv.Atoi(record[0])
		if err != nil {
			logger.Logging(fmt.Errorf("dataloader: ID converting error for %s place\n", place.Name))
			continue
		}
		place.Name = record[1]
		place.Address = record[2]
		place.Phone = record[3]
		place.Location.Longitude, err = strconv.ParseFloat(record[4], 64)
		if err != nil {
			logger.Logging(fmt.Errorf("dataloader: longitude converting error for ID %d: %s place\n", place.ID, place.Name))
			continue
		}
		place.Location.Latitude, err = strconv.ParseFloat(record[5], 64)
		if err != nil {
			logger.Logging(fmt.Errorf("dataloader: latitude converting error for ID %d: %s place\n", place.ID, place.Name))
			continue
		}

		data, err := json.Marshal(place)
		if err != nil {
			logger.Logging(fmt.Errorf("dataloader: data marshaling error for ID=%d %s place: %s\n", place.ID, place.Name, err.Error()))
			continue
		}

		err = req.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			Body:       bytes.NewReader(data),
			DocumentID: record[0],
		})
	}
	req.Close(context.Background())
	fmt.Println("Added documents:", req.Stats().NumAdded)
	fmt.Println("Failed documents:", req.Stats().NumFailed)
}

// Returns a list of items based on the limit/offset, a total number of hits and (or) an error in case of one
func (c *Client) GetPlaces(limit, offset int, index string) ([]types.Place, int, error) {
	q := map[string]interface{}{
		"size": limit,
		"from": offset,
	}
	qjson, err := json.Marshal(q)
	if err != nil {
		logger.Logging(fmt.Errorf("dataloader: cannot marshaling; %w\n", err))
		return nil, 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req := esapi.SearchRequest{
		Index:          []string{index},
		Body:           strings.NewReader(string(qjson)),
		TrackTotalHits: true,
	}

	res, err := req.Do(ctx, c.Client)
	if err != nil {
		logger.Logging(fmt.Errorf("dataloader: requset error; %w\n", err))
		return nil, 0, err
	}
	defer res.Body.Close()

	var esRes types.EsResponse
	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		logger.Logging(fmt.Errorf("dataloader: decode response error; %w\n", err))
		return nil, 0, err
	}

	places := make([]types.Place, 0, len(esRes.Hits.Hits))
	for _, hit := range esRes.Hits.Hits {
		var place types.Place
		if err := json.Unmarshal(hit.Source, &place); err != nil {
			logger.Logging(fmt.Errorf("dataloader: unmarshaling error; %w\n", err))
			continue
		}

		places = append(places, place)
	}

	return places, esRes.Hits.Total.Value, nil
}

// Returns the list of the closest places of set lat and lon in amount of (var: closest) and (or) an error in case of one
func (c *Client) GetClosestPlaces(lat, lon float64, index string) ([]types.Place, error) {
	q := fmt.Sprintf(`
	{
	"size": %d,
	"sort": [
    {
      "_geo_distance": {
        "location": {
          "lat": %v,
          "lon": %v
        },
        "order": "asc",
        "unit": "km",
        "mode": "min",
        "distance_type": "arc",
        "ignore_unmapped": true
      }
    }
	]
	}
	`,
		closest, lat, lon)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(q),
	}

	res, err := req.Do(ctx, c.Client)
	if err != nil {
		logger.Logging(fmt.Errorf("dataloader: requset error: %w\n", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Logging(fmt.Errorf("dataloader: ES requset error: %s\n", res.String()))
		return nil, fmt.Errorf("dataloader: ES request error: %s", res.String())
	}

	var esRes types.EsResponse
	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		logger.Logging(fmt.Errorf("dataloader: decode response error: %w\n", err))
		return nil, err
	}

	places := make([]types.Place, 0, len(esRes.Hits.Hits))
	for _, hit := range esRes.Hits.Hits {
		var place types.Place
		if err := json.Unmarshal(hit.Source, &place); err != nil {
			logger.Logging(fmt.Errorf("dataloader: unmarshaling error: %w\n", err))
			continue
		}

		places = append(places, place)
	}

	return places, nil
}
