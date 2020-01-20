package actions

import "github.com/gobuffalo/buffalo"

// ClassClass default implementation.
func ClassClass(c buffalo.Context) error {
	return c.Render(200, r.HTML("class/class.html"))
}

