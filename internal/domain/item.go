package domain

import "github.com/google/uuid"

type Item struct {
	ItemUID     uuid.UUID
	ChrtID      int
	TrackNumber string
	RID         string
	Name        string
	Size        string
	NmID        int
	Brand       string
	Status      int
}
