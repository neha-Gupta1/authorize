package main

import (
	"authorize/controller"
	"authorize/models"

)

func main() {

	// Connect to database
	models.ConnectDatabase()

	// Run the server
	controller.SetupServer().Run()
}
