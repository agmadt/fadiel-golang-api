package main

import (
	"golang-api/app"
	"golang-api/routes"
)

func main() {
	app.InitDatabase()
	routes.Init()
}
