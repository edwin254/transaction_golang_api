package model

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Transaction struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (trx *Transaction) BeforeCreate() (err error) {
	trx.ID = uuid.New().String()
	return
}
