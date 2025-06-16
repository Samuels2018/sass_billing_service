package models

import "time"

type Invoice struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Status        string    `json:"status"` // "pending", "paid", "cancelled"
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateInvoiceRequest struct {
	UserID        int     `json:"user_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
}
