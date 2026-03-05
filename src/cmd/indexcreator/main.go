package main

import (
	app "APIforElasticBD/pkg"
	"flag"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 3 {
		fmt.Println("indexcreator: usage $ utility -i index_name_if_it_needs")
		fmt.Println("indexcreator: by default index name is 'places'. Usage $ utility")
		return
	}

	idxname := flag.String("i", "places", "a name of the index")
	flag.Parse()

	app.ESCreateIndices(*idxname)
}
