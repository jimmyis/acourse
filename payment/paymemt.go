package payment

import (
	"context"
	"time"
)

// Payment domain
type Payment struct {
	ID            string
	User          User
	Course        Course
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        Status
	CreatedAt     time.Time
	UpdatedAt     time.Time
	At            time.Time
}

// Status is the payment status
type Status int

// Payment Status values
const (
	Pending Status = iota
	Accepted
	Rejected
	Refunded
)

// Repository is the payment repository
type Repository interface {
	Store(ctx context.Context, payment *Payment) error
	FindID(ctx context.Context, id string) (*Payment, error)
	List(ctx context.Context, status []Status, limit, offset int64) ([]*Payment, error)
	Count(ctx context.Context, status []Status) (int64, error)
}
