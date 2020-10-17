package structs

type ProductVariantOption struct {
	ID               string `db:"id"`
	ProductVariantID string `db:"product_variant_id"`
	Name             string `db:"name"`
}

type ProductVariantOptionRequest struct {
	Name string `json:"name" validate:"required"`
}

type ProductVariantOptionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
