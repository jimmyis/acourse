package appctx

import (
	"context"

	"github.com/acoshift/acourse/user"
)

type (
	userKey      struct{}
	xsrfKey      struct{}
	courseURLKey struct{}
)

// WithUser creates new context with user value
func WithUser(ctx context.Context, user *user.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) *user.User {
	x, _ := ctx.Value(userKey{}).(*user.User)
	return x
}

// WithXSRFToken creates new context with xsrf value
func WithXSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, xsrfKey{}, token)
}

// GetXSRFToken gets xsrf token from context
func GetXSRFToken(ctx context.Context) string {
	x, _ := ctx.Value(xsrfKey{}).(string)
	return x
}

// WithCourseURL creates new context with course url value
func WithCourseURL(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, courseURLKey{}, v)
}

// GetCourseURL gets course url from context
func GetCourseURL(ctx context.Context) string {
	x, _ := ctx.Value(courseURLKey{}).(string)
	return x
}
