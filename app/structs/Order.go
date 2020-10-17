package structs

import (
	helpers "golang-api/app/helpers"
	"time"
)

type Order struct {
	ID         string             `db:"id"`
	BuyerName  string             `db:"buyer_name"`
	BuyerEmail string             `db:"buyer_email"`
	Total      float32            `db:"total"`
	Message    helpers.NullString `db:"message" json:"string"`
	CreatedAt  time.Time          `db:"created_at"`
}

type OrderResponse struct {
	ID         string             `json:"id"`
	BuyerName  string             `json:"buyer_name"`
	BuyerEmail string             `json:"buyer_email"`
	Total      float32            `json:"total"`
	Message    helpers.NullString `json:"message"`
	CreatedAt  int64              `json:"created_at"`
}

type OrderPagination struct {
	Orders []OrderResponse `json:"orders"`
	Limit  int             `json:"limit"`
	Page   int             `json:"page"`
	Total  int             `json:"total"`
}

type OrderRequest struct {
	BuyerName  string                `json:"buyer_name" validate:"required"`
	BuyerEmail string                `json:"buyer_email" validate:"required,email"`
	Total      float32               `json:"total"`
	Message    helpers.NullString    `json:"message"`
	Products   []OrderProductRequest `json:"products" validate:"required,min=1,dive"`
}

type OrderProductRequest struct {
	ID       string                       `json:"id" validate:"required"`
	Quantity int                          `json:"quantity" validate:"required,numeric"`
	Variants []OrderProductVariantRequest `json:"variants" validate:"required,dive"`
}

type OrderProductVariantRequest struct {
	VariantID string `json:"variant_id" validate:"required"`
	OptionID  string `json:"option_id" validate:"required"`
}

func (order *Order) Response() OrderResponse {
	response := OrderResponse{
		ID:         order.ID,
		BuyerName:  order.BuyerName,
		BuyerEmail: order.BuyerEmail,
		Total:      order.Total,
		Message:    order.Message,
		CreatedAt:  order.CreatedAt.Unix(),
	}

	return response
}
