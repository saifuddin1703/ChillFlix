package main

import (
	"chillfix/api"
	"chillfix/config"
	"fmt"
	"log"
	"strconv"
)

func main() {
	fmt.Println("Welcome to the Video Platform!")

	// Load config
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(config.GetServerPort())
	if err != nil {
		log.Fatal(err)
	}

	// Start the api
	err = api.NewAPI(port)
	if err != nil {
		log.Fatal(err)
	}
}
