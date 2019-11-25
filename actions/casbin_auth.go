package actions

import (
	"errors"
	"github.com/casbin/casbin"
	"github.com/gobuffalo/buffalo"
	"net/http"
)

var (
	// ErrUnauthorized is returned when the user is not allowed to perform a certain action
	ErrUnauthorized = errors.New("you are unauthorized to perform the requested action")
)

// This File contains the code which makes all the Authentication stuff
func NewRBACCheckMiddleware(e *casbin.CachedEnforcer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// get user session

			return c.Error(http.StatusUnauthorized, ErrUnauthorized)
		}
	}
}

// getRole returns all roles which the user is assigned to.
func getRole(c buffalo.Context) ([]string, error) {
	if c.Value("user") != nil {
		// check Session ID first

	}

	return nil, nil
}
