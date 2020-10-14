package models

import (
	"encoding/json"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FindAllOrder(c *gin.Context) (structs.OrdersResponse, error) {

	var orders []structs.Order
	var totalResponse int
	var orderResponseSlice []structs.OrderResponse
	var ordersResponse structs.OrdersResponse
	var limit string = "10"
	var page string = "1"
	var limitInt int
	var pageInt int
	var offset string = "1"
	var sqlQuery string

	if c.Query("limit") != "" {
		limit = c.Query("limit")
	}

	limitInt, _ = strconv.Atoi(limit)
	pageInt, _ = strconv.Atoi(page)

	if c.Query("page") != "" {
		page = c.Query("page")
		pageInt, _ := strconv.ParseInt(page, 10, 64)

		offsetInt := ((int(pageInt)) - 1) * limitInt
		offset = strconv.Itoa(int(offsetInt))
	}

	db := app.GetDB()

	sqlQuery = "SELECT id, buyer_name, buyer_email, total, message, created_at FROM orders ORDER BY created_at DESC"
	sqlQuery += " LIMIT " + offset + ", " + limit

	fmt.Println(sqlQuery)

	_, err := db.Select(&orders, sqlQuery)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Error while getting all records",
		})

		return ordersResponse, err
	}

	err = db.SelectOne(&totalResponse, "SELECT count(id) as total FROM orders")

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Error while getting all records",
		})

		return ordersResponse, err
	}

	for _, order := range orders {
		orderResponseSlice = append(orderResponseSlice, order.Response())
	}

	ordersResponse = structs.OrdersResponse{
		Orders: orderResponseSlice,
		Total:  totalResponse,
		Limit:  limitInt,
		Page:   pageInt,
	}

	return ordersResponse, err
}

func FindOrder(c *gin.Context, pk string) (structs.Order, error) {

	var order structs.Order

	db := app.GetDB()

	err := db.SelectOne(&order, "SELECT id, buyer_name, buyer_email, total, message, created_at FROM orders WHERE id=?", pk)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Order not found",
		})

		return order, err
	}

	return order, err
}

func FindOrderProducts(c *gin.Context, order structs.Order) ([]structs.OrderProductProduct, error) {

	var orderProducts = []structs.OrderProductProduct{}
	var productsJSON []string

	db := app.GetDB()

	_, err := db.Select(&productsJSON, "SELECT product FROM order_products WHERE order_id=?", order.ID)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Error while getting products",
		})
		fmt.Println(err)

		return orderProducts, err
	}

	if len(productsJSON) > 0 {
		for _, product := range productsJSON {
			fmt.Println(product)
			var orderProduct = structs.OrderProductProduct{}
			productByte := []byte(product)

			err = json.Unmarshal(productByte, &orderProduct)
			if err != nil {
				c.JSON(404, gin.H{
					"message": "Error while getting products",
				})
				fmt.Println(err)

				return orderProducts, err
			}

			orderProducts = append(orderProducts, orderProduct)
		}
	}

	return orderProducts, err
}

func StoreOrder(c *gin.Context, request structs.OrderRequest) (structs.Order, error) {

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
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return order, err
	}

	return order, err
}
