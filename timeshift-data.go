package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var localTimeshiftData timeshiftData = timeshiftData{"sqlite3", "./timeshifts.db"}

type timeshiftData struct {
	dbType     string
	dbFilepath string
}

func (data timeshiftData) initDB() error {

	db, err := sql.Open(data.dbType, data.dbFilepath)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS namespaces {
			id INTEGER PRIMARY KEY,
			name TEXT
		};
		CREATE TABLE IF NOT EXISTS projects {
			id INTEGER PRIMARY KEY,
			name TEXT,
			FOREIGN KEY (namespace) REFERENCES namespaces(id),
		};
		CREATE TABLE IF NOT EXISTS timeshifts ( 
			FOREIGN KEY (project) REFERENCES projects(id),
			clock_in_time INTEGER,
			clock_out_time INTEGER
		);
	`)
	if err != nil {
		return err
	}

	res, err := stmt.exec()
	if err != nil {
		return err
	}

	return nil
}

// register the beginning of a timeshift
func (data timeshiftData) clockIn(shift timeshift) error {
	db, err := sql.Open(data.dbType, data.dbFilepath)
	if err != nil {
		return err
	}
	// FIXME: update if exists, insert otherwise
	stmt, err := db.Prepare("INSERT INTO timeshifts(project, clock_in_time), project SELECT (?,?,?,?)")
	return nil
}

// register the end of a timeshift
func (t timeshiftData) clockOut(shift timeshift) error {

	return nil
}

func (t timeshiftData) editTimeshift(targetProject project, newClockInTime time.Time, newClockOutTime time.Time) error {
	// TODO: implement me
	return nil
}

func (t timeshiftData) editProject(targetProject, newProject project) error {
	// TODO: implement me
	return nil
}

// returns all timeshifts associated with the passed project
func (t timeshiftData) getShifts(query timeshiftQuery) []timeshift {
	var timeshifts []timeshift
	return timeshifts
}

type project struct {
	name      string
	namespace string
}

type timeshift struct {
	project      project
	clockInTime  time.Time
	clockOutTime time.Time
}

type timeshiftQuery struct {
	projectName string
	namespace   string
	from        time.Time
	to          time.Time
}
