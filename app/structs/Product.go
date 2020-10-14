package structs

import (
	helpers "golang-api/app/helpers"
	"time"
)

type Product struct {
	ID          string             `db:"id" json:"id"`
	Name        string             `db:"name" json:"name"`
	Price       float32            `db:"price" json:"price"`
	Description helpers.NullString `db:"description" json:"description"`
	CreatedAt   time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at" json:"updated_at"`
}
