package controllers

import "github.com/gin-gonic/gin"

type CategoryController struct{}

func (controller CategoryController) Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "naga",
	})
}
