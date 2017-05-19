package main

import (
	"time"
)

type timeshiftDataStore struct {
	location string
	url      string
}

// register the beginning of a timeshift
func (t *timeshiftDataStore) clockIn(shift timeshift) error {

	return nil
}

// register the end of a timeshift
func (t *timeshiftDataStore) clockIn(shift timeshift) error {

	return nil
}

func (t *timeshiftDataStore) editTimeslot(targetProject project, newClockOnTime time.Time, newClockOffTime time.Time) error {
	// TODO: implement me
	return nil
}

func (t *timeshiftDataStor) editProject(targetProject, newProject project) error {
	// TODO: implement me
	return nil
}

// returns all timeshifts associated with the passed project
func (t *timeshiftDataStore) getShifts(targetProject project) []timeshift {

}

type project struct {
	name      string
	namespace string
}

type timeshift struct {
	project      project
	clockOnTime  time.Time
	clockOffTime time.Time
}
