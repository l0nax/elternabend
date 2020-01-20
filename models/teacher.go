package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

type Teacher struct {
	ID        uuid.UUID `json:"id" db:"id" rw:"r"`
	UserID    uuid.UUID `json:"-" db:"uid"`
	Name      string    `json:"name" db:"name" rw:"w"`
	Room      string    `json:"room" db:"room"` // Room describes where parents can find the teacher
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Mail      string    `json:"mail" db:"-"` // This field is only to unmarshal requests. Mail is saved in 'users' table
}

// String is not required by pop and may be deleted
func (t Teacher) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Teachers is not required by pop and may be deleted
type Teachers []Teacher

// String is not required by pop and may be deleted
func (t Teachers) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Teacher) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Teacher) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	// var err error
	// return validate.Validate(
	//         &validators.StringIsPresent{Field: t.Name, Name: "Name"},
	// ), err
	return validate.NewErrors(), nil
}
