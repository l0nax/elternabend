package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// SubjectID int       `json:"subject_id" db:"subject_id"`
type SubjectClass struct {
	ID          uuid.UUID `json:"id" db:"id"`
	SubjectName string    `json:"subject_name" db:"subject_name"`
	ClassID     uuid.UUID `json:"class_id" db:"class_id"`
	TeacherID   uuid.UUID `json:"teacher_id" db:"teacher_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (s SubjectClass) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

func (s *SubjectClass) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		// &validators.IntIsPresent{Field: s.SubjectID, Name: "SubjectID"},
		&validators.StringIsPresent{Field: s.SubjectName, Name: "SubjectName"},
		&validators.UUIDIsPresent{Field: s.ClassID, Name: "ClassID"},
		&validators.UUIDIsPresent{Field: s.TeacherID, Name: "TeacherID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *SubjectClass) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *SubjectClass) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// Exists checks if an Subject->Class entry already exists.
func (s SubjectClass) Exists(tx *pop.Connection) (bool, error) {
	err := tx.Where("subject_name = ? and class_id = ? and teacher_id = ?",
		s.SubjectName, s.ClassID, s.TeacherID)
	if err != nil {
		return false, nil
	}

	return true, nil
}
