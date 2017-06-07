package timeshifts

import (
	"errors"
	"fmt"
	"time"
)

type Project struct {
	Name      string
	Namespace string
}

type DB struct {
	DbDriver   string
	DbFilepath string
}

type Timeshift struct {
	ID           int64
	Project      Project
	ClockOnTime  time.Time
	ClockOffTime time.Time
}

type TimeshiftQuery struct {
	ProjectName string
	Namespace   string
	From        time.Time
	To          time.Time
}

type ErrTimeshiftAlreadyRunning Timeshift

func (e ErrTimeshiftAlreadyRunning) Error() string {
	originalShift := Timeshift(e)
	outputStr := "There is already a running timeshift for project %v: previous timeshift started at %v.\n"
	return fmt.Sprintf(outputStr, FormatProject(originalShift.Project), originalShift.ClockOnTime)
}

var ErrNoTimeshifts = errors.New("No timeshifts found.")
