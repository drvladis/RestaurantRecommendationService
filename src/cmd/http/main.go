package main

import (
	dl "APIforElasticBD/internal/db/dataloader"
	hl "APIforElasticBD/internal/handlers"
	"APIforElasticBD/internal/types"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	key, ok := os.LookupEnv("ES_API_KEY")
	if !ok || key == "" {
		log.Fatal("ES_API_KEY environment variable is not set")
	}
	secretKey := []byte(key)

	args := os.Args
	if len(args) > 3 {
		fmt.Println("Usage: $utility -i index")
		fmt.Println("By default using 'places' index. Usage: $utility")
		return
	}
	i := flag.String("i", "places", "a name of the used index")
	flag.Parse()
	index := *i

	var store types.Store
	store, _ = dl.Init()

	h := hl.InitHandlerParams(store, index, secretKey)

	fmt.Println("handlers are in progress.")
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Handler)
	mux.HandleFunc("/api/places", h.HandlerApi)
	mux.HandleFunc("/api/get_token", h.GetToken)
	mux.HandleFunc("/api/recommend", h.MiddleWareJWT(h.HandlerApiRecommend))
	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		log.Fatal(err)
	}
}
