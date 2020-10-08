package controllers

import "github.com/gin-gonic/gin"

type TestController struct{}

func (controller TestController) Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Test success",
	})
}
