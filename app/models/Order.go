package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func PaginateOrder(paginationParam structs.PaginationParameters) (structs.OrderPagination, error) {

	var orders []structs.Order
	var totalResponse int
	var orderResponses = []structs.OrderResponse{}
	var orderPagination structs.OrderPagination
	var offset string = "0"
	var sqlQuery string
	var limit string = strconv.Itoa(paginationParam.Limit)

	offset = strconv.Itoa(((paginationParam.Page) - 1) * paginationParam.Limit)

	db := app.GetDB()

	sqlQuery = "SELECT id, buyer_name, buyer_email, total, message, created_at FROM orders ORDER BY created_at DESC"
	sqlQuery += " LIMIT " + offset + ", " + limit

	_, err := db.Select(&orders, sqlQuery)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Query to get all order error", err)
		return orderPagination, err
	}

	err = db.SelectOne(&totalResponse, "SELECT count(id) as total FROM orders")
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Query to count all order error", err)
		return orderPagination, err
	}

	for _, order := range orders {
		orderResponses = append(orderResponses, order.Response())
	}

	orderPagination = structs.OrderPagination{
		Orders: orderResponses,
		Total:  totalResponse,
		Limit:  paginationParam.Limit,
		Page:   paginationParam.Page,
	}

	return orderPagination, err
}

func FindOrder(order structs.Order) (structs.Order, error) {

	db := app.GetDB()

	err := db.SelectOne(&order, "SELECT id, buyer_name, buyer_email, total, message, created_at FROM orders WHERE id=?", order.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Find order error", err)
		}
		return order, err
	}

	return order, err
}

func FindOrderProducts(order structs.Order) ([]structs.OrderProductProduct, error) {

	var orderProducts = []structs.OrderProductProduct{}
	var productsJSON []string

	db := app.GetDB()

	_, err := db.Select(&productsJSON, "SELECT product FROM order_products WHERE order_id=?", order.ID)
	if err != nil {
		fmt.Println("Find order products error", err)
		return orderProducts, err
	}

	if len(productsJSON) > 0 {
		for _, product := range productsJSON {
			fmt.Println(product)
			var orderProduct = structs.OrderProductProduct{}
			productByte := []byte(product)

			err = json.Unmarshal(productByte, &orderProduct)
			if err != nil {
				fmt.Println("Failed to unmarshal product from order products", err)
				return orderProducts, err
			}

			orderProducts = append(orderProducts, orderProduct)
		}
	}

	return orderProducts, err
}

func StoreOrder(request structs.OrderRequest) (structs.Order, error) {

	db := app.GetDB()

	var order = structs.Order{
		ID:         uuid.New().String(),
		BuyerName:  request.BuyerName,
		BuyerEmail: request.BuyerEmail,
		Total:      request.Total,
		Message:    request.Message,
		CreatedAt:  time.Now(),
	}

	_, err := db.Exec("INSERT INTO orders(id, buyer_name, buyer_email, total, message, created_at) VALUES (?,?,?,?,?,?)", order.ID, order.BuyerName, order.BuyerEmail, order.Total, order.Message, order.CreatedAt)
	if err != nil {
		fmt.Println("Store order error", err)
		return order, err
	}

	return order, err
}
