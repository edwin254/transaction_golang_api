package model

import (
    "time"
    "github.com/google/uuid"
)

type Transaction struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    Amount    float64   `json:"amount"`
    Currency  string    `json:"currency"`
    Sender    string    `json:"sender"`
    Receiver  string    `json:"receiver"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (t *Transaction) BeforeCreate(tx any) (err error) {
    if t.ID == "" {
        t.ID = uuid.New().String()
    }
    return nil
}
