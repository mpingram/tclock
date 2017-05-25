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

// Call this before accessing the timeshift data.
func (data timeshiftsDAO) init() error {
	// create tables if not exist
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
	// DEBUG
	rows, err := db.Query("Select * from projects")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var projectID int
		var namespaceID sql.NullInt64
		var name string
		err = rows.Scan(&projectID, &namespaceID, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println(projectID, namespaceID, name)

	}
	// END DEBUG
	projectExists, projectID := getProjectID(db, shift.project.name, shift.project.namespace)

	// create new project if not exists
	if projectExists == false {
		fmt.Println("Project dont exists yet")
		// reassign projectID to newly created project
		projectID, err = addProject(db, shift.project.name, shift.project.namespace)
		fmt.Println(projectID)
		if err != nil {
			return err
		}
	}
	// create new timeshift
	db.Exec("INSERT INTO timeshifts(project_id, clock_in_time) VALUES (?,?)", projectID, shift.clockInTime)
	if err != nil {
		return err
	}
	// no error
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

// helper function inserts new project into db
func addProject(db *sql.DB, name string, namespace string) (int64, error) {
	createProjectStmt, err := db.Prepare(`INSERT INTO projects(name, namespace_id, project_id) VALUES (?,?,NULL)`)
	var res sql.Result
	if err != nil {
		return 0, err
	}
	projectHasNamespace := namespace != ""
	if projectHasNamespace == true {
		namespaceExists, namespaceID := getNamespaceID(db, namespace)
		if namespaceExists == false {
			_, err = db.Exec(`INSERT INTO namespaces(name, namespace_id) VALUES (?, NULL)`, namespace)
			if err != nil {
				return 0, err
			}
			namespaceExists, namespaceID = getNamespaceID(db, namespace)
			if namespaceExists == false {
				panic("Namespace write failed")
			}
			fmt.Printf("%v : %v", namespace, namespaceID)
		}
		res, err = createProjectStmt.Exec(name, namespaceID)
	} else {
		res, err = createProjectStmt.Exec(name, nil)
	}
	if err != nil {
		return 0, err
	}
	projectID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return projectID, nil
}

// helper function finds projectID from name and namespace
func getProjectID(db *sql.DB, name string, namespace string) (bool, int64) {
	var exists bool
	var id int64
	hasNamespace := namespace != ""
	var err error
	fmt.Printf("namespace: %v\n", namespace)
	if hasNamespace == false {
		fmt.Println("has no namespace!")
		queryString := `
			SELECT project_id FROM projects
			WHERE projects.name=? AND projects.namespace_id IS NULL`
		err = db.QueryRow(queryString, name).Scan(&id)
	} else {
		queryString := `
				SELECT project_id FROM projects 
					INNER JOIN namespaces ON projects.namespace_id = namespaces.namespace_id 
				WHERE projects.name=? AND namespaces.name=?`
		err = db.QueryRow(queryString, name, namespace).Scan(&id)
	}
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found!")
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
