package actions

import "github.com/gobuffalo/buffalo"

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

// RoutesHandler renders the 'routes' Template
// this route is only for Developing
func RoutesHandler(c buffalo.Context) error {
	// check if we are in development Mode
	if ENV != "development" {
		return c.Render(200, r.HTML("routes.html"))
	} else {
		c.Set("isError", true)
		return c.Render(500, r.HTML("null.html", "_500.html"))
	}
}
