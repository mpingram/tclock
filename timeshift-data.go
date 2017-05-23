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
	defer db.Close()
	if err != nil {
		return err
	}
	if db == nil {
		panic("db nil")
	}
	if dbPingErr != nil {
		panic(dbPingErr)
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
	projectExists, projectID := getProjectID(db, shift.project.name, shift.project.namespace)
	// create new project if not exists
	if projectExists == false {
		fmt.Println("Project dont exists yet")
		createProjectStmt, err := db.Prepare(`INSERT INTO projects(name, namespace_id, project_id) VALUES (?,?,NULL)`)
		if err != nil {
			return err
		}
		projectHasNamespace := shift.project.namespace != ""
		if projectHasNamespace == true {
			namespaceExists, namespaceID := getNamespaceID(db, shift.project.namespace)
			if namespaceExists == false {
				_, err = db.Exec(`INSERT INTO namespaces(name, namespace_id) VALUES (?, NULL)`, shift.project.namespace)
				if err != nil {
					return err
				}
				namespaceExists, namespaceID = getNamespaceID(db, shift.project.namespace)
				if namespaceExists == false {
					panic("Namespace write failed")
				}
				fmt.Printf("%v : %v", shift.project.namespace, namespaceID)
			}
			_, err = createProjectStmt.Exec(shift.project.name, namespaceID)
		} else {
			_, err = createProjectStmt.Exec(shift.project.name, nil)
		}
		if err != nil {
			return err
		}
		// reassign projectID to the newly created projectID
		_, projectID = getProjectID(db, shift.project.name, shift.project.namespace)
	}
	// create new timeshift
	createTimeshiftStmt, err := db.Prepare("INSERT INTO timeshifts(project_id, clock_in_time) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = createTimeshiftStmt.Exec(projectID, shift.clockInTime)
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
func getProjectID(db *sql.DB, name string, namespace string) (bool, int) {
	var exists bool
	var id int
	hasNamespace := namespace != ""
	queryString := `
		  SELECT project_id FROM projects 
		    INNER JOIN namespaces ON projects.namespace_id = namespaces.namespace_id 
		  WHERE projects.name=? AND namespaces.name=?`
	var err error
	if hasNamespace == true {
		err = db.QueryRow(queryString, name, nil).Scan(&id)
	} else {
		err = db.QueryRow(queryString, name, namespace).Scan(&id)
	}
	switch {
	case err == sql.ErrNoRows:
		exists = false
	case err != nil:
		panic(err)
	default:
		exists = true
	}
	return exists, id
}

// helper function finds namespaceId from namespace
func getNamespaceID(db *sql.DB, namespace string) (bool, int) {
	var exists bool
	var id int
	queryString := `
    SELECT namespace_id from namespaces
		WHERE name=? 
		LIMIT 1"`
	err := db.QueryRow(queryString, namespace).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		exists = false
	case err != nil:
		panic(err)
	default:
		exists = true
	}
	return exists, id
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
