package routes

import (
	"golang-api/app/controllers"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	testController := new(controllers.TestController)
	categoryController := new(controllers.CategoryController)
	orderController := new(controllers.OrderController)

	router.GET("/test", testController.Test)

	router.GET("/categories", categoryController.Index)
	router.GET("/categories/:id", categoryController.Show)
	router.POST("/categories", categoryController.Store)
	router.PATCH("/categories/:id", categoryController.Update)
	router.DELETE("/categories/:id", categoryController.Delete)

	router.GET("/orders", orderController.Index)
	router.GET("/orders/:id", orderController.Show)
	router.POST("/orders", orderController.Store)

	router.Run()
}
