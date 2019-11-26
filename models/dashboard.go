package models

import (
	"time"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/pop"
	suuid "github.com/google/uuid"
)

type DashboardNew struct {
	ID         suuid.UUID   `json:"-" db:"id" rw:"r"`
	StartDate  time.Time    `json:"start_date" db:"start_date"`
	EndDate    time.Time    `json:"end_date" db:"end_date"`
	DeleteTime time.Time    `json:"delete_time" db:"delete_time"`
	ImportData binding.File `db:"-" form:"input-file"`
}

func (d *DashboardNew) ProcessImport(tx *pop.Connection) error {
	// TODO: Implement process code

	return nil
}
