package main

import (
	app "APIforElasticBD/pkg"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	arguments := os.Args
	if len(arguments) > 5 {
		fmt.Println("too many arguments. Usage: ~$ ./loadingData -i index -f filepath.csv")
		return
	}

	file := flag.String("f", "../../dataset/data.csv", "track a filpath to read")
	index := flag.String("i", "places", "a name of the index")
	flag.Parse()

	if filepath.Ext(*file) != ".csv" && *file != "" {
		fmt.Printf("file extension must be .csv: %s\n", *file)
		return
	}

	filename := filepath.Clean(*file)
	app.ESDataLoad(filename, *index)
}
