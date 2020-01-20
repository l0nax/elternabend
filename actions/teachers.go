package actions

import (
	"fmt"
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/l0nax/elternabend/internal"
	"github.com/l0nax/elternabend/models"
	"github.com/pkg/errors"
)

// TeacherList gets all Teachers.
func TeacherList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	teachers := &models.Teachers{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Teachers from the DB
	if err := q.All(teachers); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, teachers))
}

// TeacherShow gets the data for one Teacher.
func TeacherShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	// To find the Teacher the parameter teacher_id is used.
	if err := tx.Find(teacher, c.Param("teacher_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, teacher))
}

// Create adds a Teacher to the DB.
func TeacherCreate(c buffalo.Context) error {
	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	// Bind teacher to the html form elements
	if err := BindJSON(teacher, c); err != nil {
		return errors.Wrap(err, "Error while binding Params into Model")
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// we have first to create an new User which we can then connect
	// with the new Teacher
	tUser := models.User{
		Username:     teacher.Mail,
		Email:        teacher.Mail,
		ClassTeacher: true,
		Roles:        "class_teacher",
		Password:     RandomString(internal.TEACHER_PASSWORD_LENGTH),
	}

	// validate and create new User
	verrs, err := tUser.Create(tx)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "Error while validating and creating User (teacher)"))
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// response with Error
		return c.Error(422, verrs)
	}

	// add UserID
	teacher.UserID = tUser.ID

	// Validate the data from the html form
	verrs, err = tx.ValidateAndCreate(teacher)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "Error while validating and creating DB entry"))
	}

	if verrs.HasAny() {
		// delete the new created User
		err = tUser.Destroy(tx)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Error while destroying (temp) Teacher-User data"))
		}

		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// response with Error
		return c.Error(422, verrs)
	}

	// If there are no errors set a success message
	c.Set("msg", T.Translate(c, "teacher.created.success.success"))
	return c.Render(200, r.JSON(teacher))
}

// Edit renders a edit form for a Teacher.
func TeacherEdit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	if err := tx.Find(teacher, c.Param("teacher_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, teacher))
}

// Update changes a Teacher in the DB.
func TeacherUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	if err := tx.Find(teacher, c.Param("teacher_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Teacher to the html form elements
	if err := c.Bind(teacher); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(teacher)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, teacher))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "teacher.updated.success"))
	// and redirect to the teachers index page
	return c.Render(200, r.Auto(c, teacher))
}

// Destroy deletes a Teacher from the DB.
func TeacherDetroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	// To find the Teacher the parameter teacher_id is used.
	if err := tx.Find(teacher, c.Param("teacher_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(teacher); err != nil {
		return err
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", T.Translate(c, "teacher.destroyed.success"))
	// Redirect to the teachers index page
	return c.Render(200, r.Auto(c, teacher))
}

// =====================================
//	Package internal functions
// =====================================
