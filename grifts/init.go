package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/l0nax/elternabend/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
