package routes

import (
	"golang-api/app/controllers"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	testController := new(controllers.TestController)
	orderController := new(controllers.OrderController)

	router.GET("/test", testController.Test)

	router.GET("/orders", orderController.Index)
	router.GET("/orders/:id", orderController.Show)
	router.POST("/orders", orderController.Store)

	router.Run()
}
