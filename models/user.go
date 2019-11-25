package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	suuid "github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	// Salt         string    `json:"salt" db:"salt"`
	Admin        bool       `json:"-" db:"admin"`         // check if only rw: r
	ClassTeacher bool       `json:"-" db:"class_teacher"` // check if only rw: r
	TeacherID    int        `json:"-" db:"teacher_id"`    // check if only rw: r
	CreatedAt    time.Time  `json:"-" db:"created_at"`    // check if only rw: r
	UpdatedAt    time.Time  `json:"-" db:"updated_at"`    // check if only rw: r
	Roles        []string   `json:"-" db:"roles"`         // a comma seperated list of all roles
	SessionID    suuid.UUID `json:"-" db:"session_uuid"`
	Password     string     `json:"-" db:"-"`
	Email        string     `json:"email" db:"email"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Create creates a new User and adds it to the Database
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	// If the mail is not providet than use the Username as mail-address
	if len(u.Email) == 0 {
		u.Email = u.Username
	}

	// check set'ed roles
	if len(u.Roles) == 0 {
		// set – as default – 'class_teacher' as role
		u.Roles = []string{"class_teacher"}
	}

	// check if Field 'admin' is true, if so add 'admin' role to roels
	if u.Admin {
		u.Roles = append(u.Roles, "admin")
	}

	// hash Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MaxCost)
	if err != nil {
		return validate.NewErrors(), errors.Wrap(err, "Error while hashing Password: ")
	}

	u.Password = string(pwdHash)
	return tx.ValidateAndCreate(u)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Username, Name: "Username"},
		&validators.StringIsPresent{Field: u.Password, Name: "PasswordHash"},
		&validators.StringIsPresent{Field: u.Roles, Name: "Roles"}
		&UsernameNotTaken{Name: "Username", Field: u.Username, tx: tx},
		&EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&UsernameNotTaken{Name: "Username", Field: u.Username, tx: tx},
		&EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ==========> Custom Validators <==========

type UsernameNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

func (v *UsernameNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("username = ?", v.Field)
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil {
		// found a user with same username
		errors.Add(validators.GenerateKey(v.Name), log.Printf("The username %s is not available.", v.Field))
	}
}

type EmailNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

// IsValid performs the validation check for unique emails
func (v *EmailNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("email = ?", v.Field)
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil {
		// found a user with the same email
		errors.Add(validators.GenerateKey(v.Name), "An account with that email already exists.")
	}
}
