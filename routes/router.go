package routes

import (
	"golang-api/app/controllers"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	testController := new(controllers.TestController)

	router.GET("/test", testController.Test)

	router.Run()
}
