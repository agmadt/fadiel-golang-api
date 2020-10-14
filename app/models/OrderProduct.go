package models

import (
	"encoding/json"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StoreOrderProduct(c *gin.Context, order structs.Order, orderProductProduct structs.OrderProductProduct) (structs.OrderProductProduct, error) {

	db := app.GetDB()

	uuid := uuid.New().String()
	marshaledProduct, err := json.Marshal(orderProductProduct)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return orderProductProduct, err
	}

	_, err = db.Exec("INSERT INTO order_products(id, order_id, product) VALUES (?,?,?)", uuid, order.ID, marshaledProduct)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return orderProductProduct, err
	}

	return orderProductProduct, err
}
