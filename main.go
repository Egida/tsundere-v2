package main

import (
	"fmt"
	"log"
	"tsundere/source"
	"tsundere/source/database"
	"tsundere/source/master"
	"tsundere/source/models"
	"tsundere/source/models/roles"
)

var (
	Models = []models.Model{
		new(roles.Roles),
	}
)

func main() {
	fmt.Printf(source.HEADER)

	// Properly load all models
	for _, model := range Models {
		model.Serve()
	}

	// Configure & parse SQLite3 database
	err := database.Configure()
	if err != nil {
		log.Fatal(err)
	}

	// Configure & start SSH server
	err = master.Configure()
	if err != nil {
		log.Fatal(err)
	}
}
