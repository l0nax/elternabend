package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// CreateEvent renders the 'Event Creation' Template
func CreateEvent(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("dashboard/new"))
}

// CreateEventPost is the backend function witch creates and manages all the
// stuff which is needed to create a new event.
func CreateEventPost(c buffalo.Context) error {
	return nil
}

// Dashboard shows statistics and Tables with Informations about the current "Event"
func Dashboard(c buffalo.Context) error {
	// TODO: Implement a nice Dashboard
	return c.Render(http.StatusOK, r.HTML("dashboard/index"))
}
