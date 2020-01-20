package actions

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/l0nax/elternabend/models"
)

// List gets all Classes. This function is mapped to the path
// GET /classes
func ClassList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	classes := &models.Classes{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Classes from the DB
	if err := q.All(classes); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, classes))
}

// Show gets the data for one Class. This function is mapped to
// the path GET /classes/{class_id}
func ClassShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Class
	class := &models.Class{}

	// To find the Class the parameter class_id is used.
	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, class))
}

// New renders the form for creating a new Class.
// This function is mapped to the path GET /classes/new
func ClassNew(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Class{}))
}

// Create adds a Class to the DB. This function is mapped to the
// path POST /classes
func ClassCreate(c buffalo.Context) error {
	// Allocate an empty Class
	class := &models.Class{}

	// Bind class to the html form elements
	if err := BindJSON(&class, c); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(class)
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
	c.Flash().Add("success", T.Translate(c, "class.created.success"))

	// and redirect to the classes index page
	return c.Render(200, r.JSON(class))
}

// Edit renders a edit form for a Class. This function is
// mapped to the path GET /classes/{class_id}/edit
func ClassEdit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Class
	class := &models.Class{}

	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, class))
}

// Update changes a Class in the DB. This function is mapped to
// the path PUT /classes/{class_id}
func ClassUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Class
	class := &models.Class{}

	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Class to the html form elements
	if err := c.Bind(class); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(class)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, class))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "class.updated.success"))
	// and redirect to the classes index page
	return c.Render(200, r.Auto(c, class))
}

// Destroy deletes a Class from the DB. This function is mapped
// to the path DELETE /classes/{class_id}
func ClassDestroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Class
	class := &models.Class{}

	// To find the Class the parameter class_id is used.
	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(class); err != nil {
		return err
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", T.Translate(c, "class.destroyed.success"))
	// Redirect to the classes index page
	return c.Render(200, r.Auto(c, class))
}

// ClassMap is only for the API request.
// It acts like an "placeholder", so we can respond and/ or create new mappings
// between class and subject.
type JClassMap struct {
	Class        *models.Class        `json:"class"`
	SubjectClass *models.SubjectClass `json:"subject"`
	Msg          string               `json:"msg,omitempty"` // Msg is only to have an placeholder where to send a msg back
}

// ClassMap maps an teacher+subject with a specific class
func ClassMap(c buffalo.Context) error {
	classMap := &JClassMap{}

	err := BindJSON(classMap, c)
	if err != nil {
		return Wrap(err, "Error while binding request to struct")
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	ok, err = classMap.Class.Exists(tx)
	if err != nil {
		return err
	}

	// check if Class already exists
	if !ok {
		// create Class
		verrs, err := tx.ValidateAndCreate(classMap.Class)

		if err != nil {
			return Wrap(err, "Error while validating and creating DB entry")
		}

		if verrs.HasAny() {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// response with Error
			return c.Error(422, verrs)
		}
	}

	/// create Subject->Class DB entry
	// check if an equal entry already exists
	ok, err = classMap.SubjectClass.Exists(tx)
	if err != nil {
		return err
	}

	if ok {
		// return error because we can not create/ overwrite an existing
		// Subject->Class link
		classMap.Msg = "Can not overwrite an existing Subject->Class link"
	}

	// insert new Subject->Class link
	verrs, err := tx.ValidateAndCreate(classMap.SubjectClass)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.JSON(classMap))
	}

	return c.Render(200, r.JSON(classMap))
}
