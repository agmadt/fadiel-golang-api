package models

import (
	"encoding/json"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/google/uuid"
)

func StoreOrderProduct(order structs.Order, orderProductProduct structs.OrderProductProduct) (structs.OrderProductProduct, error) {

	db := app.GetDB()

	uuid := uuid.New().String()
	marshaledProduct, err := json.Marshal(orderProductProduct)
	if err != nil {
		fmt.Println("Store order product json marshal error", err)
		return orderProductProduct, err
	}

	_, err = db.Exec("INSERT INTO order_products(id, order_id, product) VALUES (?,?,?)", uuid, order.ID, marshaledProduct)
	if err != nil {
		fmt.Println("Store order product error", err)
		return orderProductProduct, err
	}

	return orderProductProduct, err
}
