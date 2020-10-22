package routes

import (
	"golang-api/app/controllers"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	categoryController := new(controllers.CategoryController)
	orderController := new(controllers.OrderController)
	productController := new(controllers.ProductController)
	mediaController := new(controllers.MediaController)
	testController := new(controllers.TestController)

	router.GET("/test", testController.Test)

	router.GET("/categories", categoryController.Index)
	router.POST("/categories/:id", categoryController.Store)
	router.GET("/categories/:id", categoryController.Show)
	router.PATCH("/categories/:id", categoryController.Update)
	router.DELETE("/categories/:id", categoryController.Delete)

	router.GET("/orders", orderController.Index)
	router.POST("/orders", orderController.Store)
	router.GET("/orders/:id", orderController.Show)

	router.GET("/products", productController.Index)
	router.POST("/products", productController.Store)
	router.GET("/products/:id", productController.Show)
	router.PATCH("/products/:id", productController.Update)
	router.DELETE("/products/:id", productController.Delete)

	router.POST("/media", mediaController.Store)

	router.Run()
}
