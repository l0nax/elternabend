package actions

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gobuffalo/buffalo"
	"github.com/l0nax/elternabend/models"
	"github.com/pkg/errors"
)

var (
	// ErrUnauthorized is returned when the user is not allowed to perform a certain action
	ErrUnauthorized = errors.New("you are unauthorized to perform the requested action")
	// ErrNotAuthorized is returned when the user is not authorized
	ErrNotAuthorized = errors.New("you are not authorized, please login!")
)

// This File contains the code which makes all the Authentication stuff
func NewRBACCheckMiddleware(e *casbin.CachedEnforcer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			roles, err := getRole(c)
			if err != nil {
				return err
			}

			res, err := e.EnforceSafe(roles, c.Request().URL.Path, c.Request().Method)
			if err != nil {
				return errors.WithStack(err)
			}

			// if User is authorized to access the URI return
			if res {
				return next(c)
			}

			return c.Error(http.StatusUnauthorized, ErrUnauthorized)
		}
	}
}

// getRole returns all roles which the user is assigned to.
func getRole(c buffalo.Context) ([]string, error) {
	if user := c.Value("user").(*models.User); user != nil {
		roles := strings.Split(user.Roles, ",") // user.Roles

		// return 'anonymous' if no roles is declared
		if len(roles) == 0 {
			roles[0] = "anonymous"
		}

		return roles, nil

	}
	return nil, ErrNotAuthorized
}
