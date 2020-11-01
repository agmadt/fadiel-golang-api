package routes

import (
	"golang-api/app/controllers"
	"golang-api/app/middlewares"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.Use(middlewares.CORSMiddleware())

	authController := new(controllers.AuthController)
	categoryController := new(controllers.CategoryController)
	orderController := new(controllers.OrderController)
	productController := new(controllers.ProductController)
	mediaController := new(controllers.MediaController)
	testController := new(controllers.TestController)

	router.POST("/auth/login", authController.Login)

	router.GET("/test", testController.Test)

	router.GET("/categories", categoryController.Index)
	router.POST("/categories", middlewares.JWTAuth(), categoryController.Store)
	router.GET("/categories/:id", middlewares.JWTAuth(), categoryController.Show)
	router.PATCH("/categories/:id", middlewares.JWTAuth(), categoryController.Update)
	router.DELETE("/categories/:id", middlewares.JWTAuth(), categoryController.Delete)

	router.GET("/orders", middlewares.JWTAuth(), orderController.Index)
	router.POST("/orders", orderController.Store)
	router.GET("/orders/:id", orderController.Show)

	router.GET("/products", productController.Index)
	router.POST("/products", middlewares.JWTAuth(), productController.Store)
	router.GET("/products/:id", productController.Show)
	router.PATCH("/products/:id", middlewares.JWTAuth(), productController.Update)
	router.DELETE("/products/:id", middlewares.JWTAuth(), productController.Delete)

	router.POST("/media", middlewares.JWTAuth(), mediaController.Store)

	router.Run()
}
