package app

import (
	"time"
)

// UserView type
type UserView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}

// UserTinyView type
type UserTinyView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
}

// UserMeView type
type UserMeView struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Photo    string    `json:"photo"`
	AboutMe  string    `json:"aboutMe"`
	Role     *RoleView `json:"role"`
}

// RoleView type
type RoleView struct {
	Admin      bool `json:"admin"`
	Instructor bool `json:"instructor"`
}

// CourseView type
type CourseView struct {
	ID               string                      `json:"id"`
	CreatedAt        time.Time                   `json:"createdAt"`
	UpdatedAt        time.Time                   `json:"updatedAt"`
	Owner            *UserTinyView               `json:"owner"`
	Title            string                      `json:"title"`
	ShortDescription string                      `json:"shortDescription"`
	Description      string                      `json:"description"`
	Photo            string                      `json:"photo"`
	Start            time.Time                   `json:"start"`
	URL              string                      `json:"url"`
	Video            string                      `json:"video"`
	Type             string                      `json:"type"`
	Price            float64                     `json:"price"`
	DiscountedPrice  float64                     `json:"discountedPrice"`
	Student          int                         `json:"student"`
	Contents         CourseContentCollectionView `json:"contents"`
	Enrolled         bool                        `json:"enrolled"`
	Enroll           bool                        `json:"enroll"`
	Public           bool                        `json:"public"`
	Owned            bool                        `json:"owned"`
	Attend           bool                        `json:"attend"`
	Assignment       bool                        `json:"assignment"`
	Purchase         bool                        `json:"purchase"`
	Discount         bool                        `json:"discount"`
}

// CoursePublicView type
type CoursePublicView struct {
	ID               string        `json:"id"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	Owner            *UserTinyView `json:"owner"`
	Title            string        `json:"title"`
	ShortDescription string        `json:"shortDescription"`
	Description      string        `json:"description"`
	Photo            string        `json:"photo"`
	Start            time.Time     `json:"start"`
	URL              string        `json:"url"`
	Video            string        `json:"video"`
	Type             string        `json:"type"`
	Price            float64       `json:"price"`
	DiscountedPrice  float64       `json:"discountedPrice"`
	Student          int           `json:"student"`
	Enroll           bool          `json:"enroll"`
	Purchase         bool          `json:"purchase"`
	Discount         bool          `json:"discount"`
}

// CourseTinyView type
type CourseTinyView struct {
	ID               string        `json:"id"`
	Owner            *UserTinyView `json:"owner"`
	Title            string        `json:"title"`
	ShortDescription string        `json:"shortDescription"`
	Photo            string        `json:"photo"`
	Start            time.Time     `json:"start"`
	URL              string        `json:"url"`
	Type             string        `json:"type"`
	Price            float64       `json:"price"`
	DiscountedPrice  float64       `json:"discountedPrice"`
	Student          int           `json:"student"`
	Discount         bool          `json:"discount"`
}

// CourseMiniView type
type CourseMiniView struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// CourseTinyCollectionView type
type CourseTinyCollectionView []*CourseTinyView

// CourseContentView type
type CourseContentView struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Video       string `json:"video"`
	DownloadURL string `json:"downloadURL"`
}

// CourseContentCollectionView type
type CourseContentCollectionView []*CourseContentView

// PaymentView type
type PaymentView struct {
	ID            string          `json:"id"`
	User          *UserTinyView   `json:"user"`
	Course        *CourseMiniView `json:"course"`
	OriginalPrice float64         `json:"originalPrice"`
	Price         float64         `json:"price"`
	Code          string          `json:"code"`
	URL           string          `json:"url"`
	Status        string          `json:"status"`
}

// PaymentCollectionView type
type PaymentCollectionView []*PaymentView