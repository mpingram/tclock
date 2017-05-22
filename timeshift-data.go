package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type timeshiftsDAO struct {
	dbDriver   string
	dbFilepath string
}

func (data timeshiftsDAO) init() error {

	db, err := sql.Open(data.dbDriver, data.dbFilepath)
	dbPingErr := db.Ping()
	if dbPingErr != nil {
		panic(dbPingErr)
	}
	defer db.Close()
	if err != nil {
		return err
	}
	if db == nil {
		panic("db nil")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS namespaces (
			namespace_id INTEGER PRIMARY KEY,
			name TEXT
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			project_id INTEGER PRIMARY KEY,
			namespace_id INTEGER,
			name TEXT,
			FOREIGN KEY (namespace_id) REFERENCES namespaces(namespace_id)
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS timeshifts ( 
			id INTEGER PRIMARY KEY,
			project_id INTEGER,
			clock_in_time INTEGER,
			clock_out_time INTEGER,
			FOREIGN KEY (project_id) REFERENCES projects(project_id)
		)
	`)

	if err != nil {
		return err
	}

	return nil
}

// register the beginning of a timeshift
func (data timeshiftsDAO) clockIn(shift timeshift) error {

	db, err := sql.Open(data.dbDriver, data.dbFilepath)
	defer db.Close()
	if err != nil {
		return err
	}
	idExists, projectID := getProjectID(db, shift.project.name, shift.project.namespace)
	// DEBUG
	fmt.Println(projectID)
	// END DEBUG
	stmt, err := db.Prepare("INSERT INTO timeshifts(project_id, clock_in_time) VALUES (?,?)")
	if err != nil {
		return err
	}
	switch {
	case idExists == true:
		// DEBUG
		fmt.Println("Project do exists yet")
		_, err = stmt.Exec(projectID, shift.clockInTime)
	case idExists == false:
		fmt.Println("Project dont exists yet")
		_, err = stmt.Exec(nil, shift.clockInTime)
	}

	if err != nil {
		return err
	}

	return nil
}

// register the end of a timeshift
func (data timeshiftsDAO) clockOut(shift timeshift) error {
	// find projectID
	return nil
}

func (data timeshiftsDAO) editTimeshift(targetProject project, newClockInTime time.Time, newClockOutTime time.Time) error {
	// TODO: implement me
	return nil
}

func (data timeshiftsDAO) editProject(targetProject, newProject project) error {
	// TODO: implement me
	return nil
}

// returns all timeshifts associated with the passed project
func (data timeshiftsDAO) getShifts(query timeshiftQuery) []timeshift {
	var timeshifts []timeshift
	return timeshifts
}

// helper function finds projectID from name and namespace
func getProjectID(db *sql.DB, name string, namespace string) (idExists bool, id int) {

	queryString := `
		  SELECT project_id FROM projects 
		    INNER JOIN namespaces ON projects.namespace_id = namespaces.namespace_id 
		  WHERE projects.name=? AND namespaces.name=?`

	err := db.QueryRow(queryString, name, namespace).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false, 0
	case err != nil:
		panic(err)
	default:
		return true, id
	}
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
