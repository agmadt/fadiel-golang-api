package main

import (
	"golang-api/app"
	"golang-api/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.InitDatabase()
	routes.Init()
}
