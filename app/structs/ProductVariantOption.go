package structs

type ProductVariantOption struct {
	ID               string `db:"id" json:"id"`
	ProductVariantID string `db:"product_variant_id" json:"product_variant_id"`
	Name             string `db:"name" json:"name"`
}
