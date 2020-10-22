package structs

import "time"

type Media struct {
	ID        string    `db:"id"`
	Filename  string    `db:"filename"`
	Location  string    `db:"location"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type MediaResponse struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

type MediaRequest struct {
	Media []byte `json:"media"`
}

func (media *Media) Response() MediaResponse {
	return MediaResponse{
		ID:       media.ID,
		Location: media.Location,
	}
}
