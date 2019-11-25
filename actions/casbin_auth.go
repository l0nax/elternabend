package actions

import (
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gobuffalo/buffalo"
	"github.com/l0nax/elternabend/models"
	"github.com/pkg/errors"
)

var (
	// ErrUnauthorized is returned when the user is not allowed to perform a certain action
	ErrUnauthorized = errors.New("You are unauthorized to perform the requested action")
	// ErrNotAuthorized is returned when the user is not authorized
	ErrNotAuthorized = errors.New("you are not authorized, please login!")
)

// This File contains the code which makes all the Authentication stuff
func NewRBACCheckMiddleware(e *casbin.Enforcer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			roles, err := getRole(c)
			if err != nil {
				return err
			}

			log.Printf("Found roles: '%v'", roles)

			// remove the last slash from the URI
			uri := strings.TrimRight(c.Request().URL.Path, "/")

			for _, role := range roles {
				res, err := e.EnforceSafe(role, uri, c.Request().Method)
				if err != nil {
					return errors.Wrap(errors.WithStack(err), "Error while creating Enforcer")
				}

				// if User is authorized to access the URI return
				// and break this looph
				if res {
					return next(c)
				}

				log.Printf("Checking group '%s' against URI '%s' with Method '%s'", role, uri, c.Request().Method)
			}

			// when returning => no role found that has the right
			// to access the page
			return c.Error(http.StatusUnauthorized, ErrUnauthorized)
		}
	}
}

// getRole returns all roles which the user is assigned to.
func getRole(c buffalo.Context) ([]string, error) {
	if user, ok := c.Value("user").(*models.User); ok || user != nil {
		roles := strings.Split(user.Roles, ",") // user.Roles

		// return 'anonymous' if no roles is declared
		if len(roles) == 0 {
			goto anonymous
		}

		return roles, nil
	}

anonymous:
	return []string{"anonymous"}, nil
}
