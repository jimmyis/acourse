package user

import (
	"context"
	"errors"
	"time"
)

// User type
type User struct {
	ID        string
	Username  string
	Name      string
	Email     string
	AboutMe   string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      Role
}

// Repository is the user storage
type Repository interface {
	// Store creates/updates user to storage
	Store(ctx context.Context, user *User) error

	// FindID finds user from id
	FindID(ctx context.Context, id string) (*User, error)

	// FindEmail finds user from email
	FindEmail(ctx context.Context, email string) (*User, error)

	// FindUsername finds user from username
	FindUsername(ctx context.Context, username string) (*User, error)

	// List lists users
	List(ctx context.Context, limit, offset int64) ([]*User, error)

	// Count counts users
	Count(ctx context.Context) (int64, error)
}

// Errors
var (
	ErrNotFound = errors.New("user: not found")
)
