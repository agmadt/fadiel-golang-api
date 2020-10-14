package structs

type Image struct {
	ID        string `db:"id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	Image     string `db:"image" json:"image"`
}
