package entity

import "github.com/google/uuid"

type Topup struct {
	UserID uuid.UUID
	Amount int
}

type TopupMsg struct {
	Id    string `json:"id"`
	Topup Topup  `json:"topup"`
}
