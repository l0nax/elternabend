package actions

import "github.com/gobuffalo/buffalo"

// SubjectClassSubjectClass default implementation.
func SubjectClassSubjectClass(c buffalo.Context) error {
	return c.Render(200, r.HTML("subject_class/subject_class.html"))
}

