package course

import (
	"context"
	"os/user"
	"time"
)

// Course type
type Course struct {
	ID           string
	Option       Option
	Owner        user.User
	EnrollCount  int64
	Title        string
	ShortDesc    string
	Desc         string
	Image        string
	Start        time.Time
	URL          string
	Type         Type
	Price        float64
	Discount     float64
	Content      []*Content
	EnrollDetail string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Option is the course option
type Option struct {
	Public     bool
	Enroll     bool
	Attend     bool
	Assignment bool
	Discount   bool
}

// Type is the course type
type Type int

// Type values
const (
	_ Type = iota
	Live
	Video
	EBook
)

// Content is the course content
type Content struct {
	ID          string
	Title       string
	Desc        string
	VideoID     string
	VideoType   VideoType
	DownloadURL string
}

// VideoType is the course content video type
type VideoType int

// VideoType values
const (
	_ VideoType = iota
	Youtube
)

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
