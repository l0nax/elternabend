package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/l0nax/elternabend/models"
)

// List gets all Subjects. This function is mapped to the path
// GET /subjects
func SubjectList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	subjects := &models.Subjects{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Subjects from the DB
	if err := q.All(subjects); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, subjects))
}

// Show gets the data for one Subject. This function is mapped to
// the path GET /subjects/{subject_id}
func SubjectShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Subject
	subject := &models.Subject{}

	// To find the Subject the parameter subject_id is used.
	if err := tx.Find(subject, c.Param("subject_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, subject))
}

// New renders the form for creating a new Subject.
// This function is mapped to the path GET /subjects/new
func SubjectNew(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Subject{}))
}

// Create adds a Subject to the DB. This function is mapped to the
// path POST /subjects
func SubjectCreate(c buffalo.Context) error {
	// Allocate an empty Subject
	subject := &models.Subject{}

	// Bind subject to the html form elements
	if err := BindJSON(subject, c); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := CreateSubject(subject, tx)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Error(422, verrs)
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "subject.created.success"))

	// and redirect to the subjects index page
	return c.Render(200, r.JSON(subject))
}

// Edit renders a edit form for a Subject. This function is
// mapped to the path GET /subjects/{subject_id}/edit
func SubjectEdit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Subject
	subject := &models.Subject{}

	if err := tx.Find(subject, c.Param("subject_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, subject))
}

// Update changes a Subject in the DB. This function is mapped to
// the path PUT /subjects/{subject_id}
func SubjectUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Subject
	subject := &models.Subject{}

	if err := tx.Find(subject, c.Param("subject_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Subject to the html form elements
	if err := c.Bind(subject); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(subject)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, subject))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "subject.updated.success"))
	// and redirect to the subjects index page
	return c.Render(200, r.Auto(c, subject))
}

// Destroy deletes a Subject from the DB. This function is mapped
// to the path DELETE /subjects/{subject_id}
func SubjectDestroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Subject
	subject := &models.Subject{}

	// To find the Subject the parameter subject_id is used.
	if err := tx.Find(subject, c.Param("subject_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(subject); err != nil {
		return err
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", T.Translate(c, "subject.destroyed.success"))
	// Redirect to the subjects index page
	return c.Render(200, r.Auto(c, subject))
}

// CreateSubject creates a new Subject DB entry
func CreateSubject(s *models.Subject, tx *pop.Connection) (*validate.Errors, error) {
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(s)
	if err != nil {
		return nil, err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		return verrs, nil
	}

	return nil, nil
}

// // SubjectExists checks if an subject already exists
// func SubjectExists(s *models.Subject, tx *pop.Connection) (*validate.Errors, error) {

// }
