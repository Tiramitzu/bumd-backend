package models

import "github.com/google/uuid"

type PendidikanModel struct {
	ID   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
}
