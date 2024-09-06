package main

import (
	"github.com/anirudhchy/gobank/api"
	"github.com/anirudhchy/gobank/storage"
	"log"
)

func main() {
	// constructor that initializes the connection with db
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// init method to create the accounts table if not already exists
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", store)

	// constructor to create an instance of the API server with port 3000
	server := api.NewAPIServer(":3000", store)

	// Run method to initialize the http server on the given route and with a router
	server.Run()
}
