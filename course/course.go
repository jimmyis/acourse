package course

import (
	"context"
)

// Course type
type Course struct {
	ID     string
	Option Option
}

// Option is the course option
type Option struct {
	Public     bool
	Enroll     bool
	Attend     bool
	Assignment bool
	Discount   bool
}

// Repository is the course storage
type Repository interface {
	// Store stores course in storage
	Store(ctx context.Context, course *Course) error

	// FindID finds course by id
	FindID(ctx context.Context, id string) (*Course, error)

	// FindURL finds course by url
	FindURL(ctx context.Context, u string) (*Course, error)

	// List lists courses
	List(ctx context.Context, limit, offset int64) ([]*Course, error)

	// Count counts courses
	Count(ctx context.Context) (int64, error)
}
