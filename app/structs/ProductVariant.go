package structs

type ProductVariant struct {
	ID        string `db:"id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	Name      string `db:"name" json:"name"`
}

type ProductVariantRequest struct {
	Name    string                        `json:"name" validate:"required"`
	Options []ProductVariantOptionRequest `json:"options" validate:"required,min=1,dive"`
}

type ProductVariantResponse struct {
	ID      string                         `json:"id"`
	Name    string                         `json:"name"`
	Options []ProductVariantOptionResponse `json:"options"`
}
