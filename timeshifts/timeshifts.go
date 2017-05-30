package timeshifts

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// register the beginning of a timeshift
func (dbRef DB) ClockOn(project Project, forceOverwrite bool) (Timeshift, error) {
	db, err := sql.Open(dbRef.DbDriver, dbRef.DbFilepath)
	defer db.Close()
	if err != nil {
		return Timeshift{}, err
	}
	projectExists, projectID, err := getProjectID(db, project)
	if err != nil {
		return Timeshift{}, err
	}
	if projectExists == false {
		projectID, err = addProject(db, project)
		if err != nil {
			return Timeshift{}, err
		}
	}
	// if there's already a clocked-on timeshift for this project,
	// 	and the method wasn't instructed to force overwrite it, respond
	// 	with ErrTimeshiftAlreadyRunning
	timeshiftAlreadyRunning, prevTimeshift, err := isTimeshiftAlreadyRunning(db, project)
	if err != nil {
		return Timeshift{}, err
	} else if timeshiftAlreadyRunning == true && !forceOverwrite {
		return Timeshift{}, ErrTimeshiftAlreadyRunning(prevTimeshift)
	}
	// create new timeshift
	clockOnTime := time.Now()
	_, err = db.Exec("INSERT INTO timeshifts(project_id, clock_on_time) VALUES (?,?)", projectID, clockOnTime.Unix())
	if err != nil {
		return Timeshift{}, err
	}
	return Timeshift{Project: project, ClockOnTime: clockOnTime}, nil
}

// register the end of a timeshift
func (dbRef DB) ClockOff(project Project) (Timeshift, error) {
	return Timeshift{}, nil
}

func (dbRef DB) GetRunningShifts() ([]Timeshift, error) {
	return []Timeshift{}, nil
}

func (dbRef DB) EditTimeshift(targetProject Project, newClockOnTime time.Time, newClockOffTime time.Time) error {
	// TODO: implement me
	return nil
}

func (dbRef DB) EditProject(targetProject, newProject Project) error {
	// TODO: implement me
	return nil
}

// returns all timeshifts associated with the passed project
func (dbRef DB) GetShifts(query TimeshiftQuery) []Timeshift {
	var timeshifts []Timeshift
	return timeshifts
}
