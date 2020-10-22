package models

import (
	"database/sql"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
)

func FindUser(user structs.User) (structs.User, error) {

	db := app.GetDB()

	err := db.SelectOne(&user, "SELECT id, name, email, password FROM users WHERE email=?", user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Find user error", err)
		}
		return user, err
	}

	return user, nil
}
