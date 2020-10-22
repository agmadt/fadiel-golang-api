package models

import (
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
	"time"

	"github.com/google/uuid"
)

func StoreMedia(media structs.Media) (structs.Media, error) {

	media = structs.Media{
		ID:        uuid.New().String(),
		Filename:  media.Filename,
		Location:  media.Location,
		Type:      media.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := app.GetDB()

	_, err := db.Exec("INSERT INTO media(id, filename, location, type, created_at, updated_at) VALUES (?,?,?,?, ?, ?)", media.ID, media.Filename, media.Location, media.Type, media.CreatedAt, media.UpdatedAt)
	if err != nil {
		fmt.Println("Store media error", err)
		return media, err
	}

	return media, nil
}
