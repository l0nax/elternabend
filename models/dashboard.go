package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// also named "Event"
type Dashboard struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StartDate time.Time `json:"start_date" db:"start_date"` // From which date users can sign up.
	EndDat    time.Time `json:"end_dat" db:"end_dat"`       // From which date the sign up is closed. This describes NOT when the event is finished!
	KillDate  time.Time `json:"kill_date" db:"kill_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (d Dashboard) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Dashboards is not required by pop and may be deleted
type Dashboards []Dashboard

// String is not required by pop and may be deleted
func (d Dashboards) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *Dashboard) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *Dashboard) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *Dashboard) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
