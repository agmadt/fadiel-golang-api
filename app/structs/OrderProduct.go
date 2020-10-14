package structs

import helpers "golang-api/app/helpers"

type OrderProduct struct {
	Order
	Products []OrderProductProduct
}

type OrderProductResponse struct {
	OrderResponse
	Products []OrderProductProduct `json:"products"`
}

type OrderProductProduct struct {
	ID          string                        `db:"id" json:"id"`
	Name        string                        `db:"name" json:"name"`
	Price       float32                       `db:"price" json:"price"`
	Description helpers.NullString            `db:"description" json:"description"`
	Quantity    int                           `json:"quantity"`
	Subtotal    float32                       `json:"subtotal"`
	Images      []OrderProductProductImage    `json:"images"`
	Variants    []OrderProductProductVariant  `json:"variants"`
	Categories  []OrderProductProductCategory `json:"categories"`
}

type OrderProductProductImage struct {
	Image string `json:"image"`
}

type OrderProductProductVariant struct {
	VariantID         string `json:"variant_id"`
	VariantName       string `json:"variant_name"`
	VariantOptionID   string `json:"variant_option_id"`
	VariantOptionName string `json:"variant_option_name"`
}

type OrderProductProductCategory struct {
	Name string `json:"name"`
}
