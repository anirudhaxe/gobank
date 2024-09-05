package main

import (
	"github.com/anirudhchy/gobank/api"
	"github.com/anirudhchy/gobank/storage"
	"log"
)

func main() {
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", store)

	server := api.NewAPIServer(":3000", store)
	server.Run()
}
