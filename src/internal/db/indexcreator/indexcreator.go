package indexcreator

import (
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

var index string = "places"

// Returns error if index doesn't exists and nil otherwise
func checkIndex(client *elasticsearch.Client) error {
	res, err := client.Indices.Exists([]string{index})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("index doesn't exist")
	}

	return nil
}

// Returns error if index didn't create and nil otherwise
func createIndex(client *elasticsearch.Client) error {
	if err := checkIndex(client); err != nil {
		if err.Error() != "index doesn't exist" {
			return err
		}
	}

	mapping :=
		`
	{
	  "settings": {
	    "number_of_shards": 1,
		"max_result_window": 20000
	},
	  "mappings": {
	    "properties": {
		  "name": {
		    "type":  "text"
		  },
	      "address": {
		    "type":  "text"
		  },
		  "phone": {
			"type":  "text"
		  },
		  "location": {
		    "type": "geo_point"
		  }
		}
	  }
	}
	`

	res, err := client.Indices.Create(index, client.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("%s", res.String())
	}

	return nil
}

// Returns error if index didn't delete and nil otherwise
func deleteIndex(client *elasticsearch.Client) error {
	if err := checkIndex(client); err != nil {
		if err.Error() == "index doesn't exist" {
			return nil
		} else {
			return err
		}
	}
	res, err := client.Indices.Delete([]string{index})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("%s", res.String())
	}

	return nil
}

// Creates index with the specified name
func IndexCreator(idxname string) {
	index = idxname

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Creating client error:", err)
	}
	fmt.Println("The client was created")

	err = deleteIndex(client)
	if err != nil {
		log.Fatal("Deleting indices error:", err)
	}
	fmt.Println("Old index was successfully deleted if existed")

	err = createIndex(client)
	if err != nil {
		log.Fatal("Creating indices error:", err)
	}
	fmt.Printf("Index '%s' was successfully created and mapped!\n", index)
}
