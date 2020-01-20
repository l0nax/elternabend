package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/l0nax/elternabend/models"
	"github.com/pkg/errors"
)

// CreateEvent renders the 'Event Creation' Template
func CreateEvent(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("dashboard/new"))
}

// CreateEventPost is the backend function witch creates and manages all the
// stuff which is needed to create a new event.
func CreateEventPost(c buffalo.Context) error {
	data := &models.DashboardNew{}

	// try to bind the Form into the Model
	if err := c.Bind(data); err != nil {
		return errors.Wrap(errors.WithStack(err), "Error while binding form into Model struct")
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("There was an error while getting 'tx' from context"))
	}

	validate, err := tx.ValidateAndCreate(data)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Error while adding the data into the database")
	}

	// check if the validation has returned any error
	if validate.HasAny() {
		// make data struct and errors in html template accessible
		c.Set("data", data)
		c.Set("errors", validate)

		return c.Render(http.StatusUnprocessableEntity, r.HTML("dashboard/new"))
	}

	c.Flash().Add("success", "Successfully created new Event")

	return c.Redirect(http.StatusTemporaryRedirect, "/d/")
}

// Dashboard shows statistics and Tables with Informations about the current "Event"
func Dashboard(c buffalo.Context) error {
	// TODO: Implement a nice Dashboard
	return c.Render(http.StatusOK, r.HTML("dashboard/index"))
}

// DashboardDashboard default implementation.
func DashboardDashboard(c buffalo.Context) error {
	return c.Render(200, r.HTML("dashboard/dashboard.html"))
}

