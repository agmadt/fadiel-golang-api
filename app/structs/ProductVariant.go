package structs

type ProductVariant struct {
	ID        string `db:"id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	Name      string `db:"name" json:"name"`
}
