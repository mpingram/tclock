package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// register the beginning of a timeshift
func (data timeshiftsDAO) clockOn(shift timeshift, forceOverwrite bool) error {
	db, err := sql.Open(data.dbDriver, data.dbFilepath)
	defer db.Close()
	if err != nil {
		return err
	}
	projectExists, projectID, err := getProjectID(db, shift.project)
	if err != nil {
		return err
	}
	if projectExists == false {
		projectID, err = addProject(db, shift.project.name, shift.project.namespace)
		if err != nil {
			return err
		}
	}
	// if there's already a clocked-on timeshift for this project,
	// 	and the method wasn't instructed to force overwrite it, respond
	// 	with ErrTimeshiftAlreadyRunning
	timeshiftAlreadyRunning, prevTimeshift, err := isTimeshiftAlreadyRunning(db, shift)
	if err != nil {
		return err
	} else if timeshiftAlreadyRunning == true && !forceOverwrite {
		return ErrTimeshiftAlreadyRunning(prevTimeshift)
	}
	// create new timeshift
	clockOnTime := shift.clockOnTime.Unix()
	_, err = db.Exec("INSERT INTO timeshifts(project_id, clock_on_time) VALUES (?,?)", projectID, clockOnTime)
	if err != nil {
		return err
	}
	return nil
}

// register the end of a timeshift
func (data timeshiftsDAO) clockOff(shift timeshift) error {
	// find projectID
	return nil
}

func (data timeshiftsDAO) editTimeshift(targetProject project, newClockOnTime time.Time, newClockOffTime time.Time) error {
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

// helper function checks if there is an unclosed timeshift for shift's project
func isTimeshiftAlreadyRunning(db *sql.DB, shift timeshift) (bool, timeshift, error) {
	var prevShift timeshift
	exists, projectID, err := getProjectID(db, shift.project)
	if exists == false {
		return false, prevShift, nil
	} else if err != nil {
		return false, prevShift, err
	} else {
		// search in timeshifts for first shift matching clock off time
		// --
		row := db.QueryRow(`
			SELECT projects.name, namespaces.name, timeshifts.clock_on_time FROM timeshifts
				INNER JOIN projects ON timeshifts.project_id=projects.project_id
					LEFT JOIN namespaces ON projects.namespace_id=namespaces.namespace_id
			WHERE timeshifts.project_id=? AND timeshifts.clock_off_time IS NULL
			LIMIT 1`, projectID)
		// reconstruct a timeshift from the result
		var name string
		var maybeNamespace sql.NullString
		var clockOnTimeUnix int64
		err = row.Scan(&name, &maybeNamespace, &clockOnTimeUnix)
		switch {
		case err == sql.ErrNoRows:
			// DEBUG
			fmt.Println("Timeshift is NOT already running.")
			return false, prevShift, nil
		case err != nil:
			return false, prevShift, err
		default:
			fmt.Println("Timeshift IS already running.")
			var proj project
			if maybeNamespace.Valid == true {
				namespace := maybeNamespace.String
				proj = project{name, namespace}
			} else {
				proj = project{name: name}
			}
			clockOnTime := time.Unix(clockOnTimeUnix, 0)
			prevShift = timeshift{project: proj, clockOnTime: clockOnTime}
			return true, prevShift, nil
		}
	}
}

// helper function inserts new namespace into db
func addNamespace(db *sql.DB, namespace string) (int64, error) {
	createNamespaceStmt, err := db.Prepare(`INSERT INTO namespaces(name, namespace_id) VALUES (?, NULL)`)
	if err != nil {
		return 0, err
	}
	res, err := createNamespaceStmt.Exec(namespace)
	if err != nil {
		return 0, err
	}
	namespaceID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return namespaceID, nil
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
			// reassign namespaceID to newly created value
			namespaceID, err = addNamespace(db, namespace)
			if err != nil {
				return 0, err
			}
			fmt.Printf("%v:%v", namespace, namespaceID)
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
func getProjectID(db *sql.DB, proj project) (exists bool, id int64, err error) {
	hasNamespace := proj.namespace != ""
	if hasNamespace == false {
		fmt.Println("has no namespace!")
		queryString := `
			SELECT project_id FROM projects
			WHERE projects.name=? AND projects.namespace_id IS NULL`
		err = db.QueryRow(queryString, proj.name).Scan(&id)
	} else {
		fmt.Printf("namespace: %v\n", proj.namespace)
		queryString := `
				SELECT project_id FROM projects 
					INNER JOIN namespaces ON projects.namespace_id = namespaces.namespace_id 
				WHERE projects.name=? AND namespaces.name=?`
		err = db.QueryRow(queryString, proj.name, proj.namespace).Scan(&id)
	}
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found!")
		exists = false
		id = 0
		err = nil
	case err != nil:
		exists = false
		id = 0
	default:
		exists = true
		err = nil
	}
	return exists, id, err
}

// DEBUG
func (data timeshiftsDAO) printDB() error {
	db, err := sql.Open(data.dbDriver, data.dbFilepath)
	if err != nil {
		return err
	}
	fmt.Println("==== timeshifts ====")
	fmt.Println("timeshift_id\tproject_id\tclock_on_time\tclock_off_time")
	timeshifts, err := db.Query("SELECT * FROM timeshifts")
	if err != nil {
		panic(err)
	}
	for timeshifts.Next() {
		var timeshiftID, projectID int64
		var clockOnTimeUnix, clockOffTimeUnix int64
		var clockOnTime, clockOffTime time.Time
		timeshifts.Scan(&timeshiftID, &projectID, &clockOnTimeUnix, &clockOffTimeUnix)
		clockOnTime = time.Unix(clockOnTimeUnix, 0)
		clockOffTime = time.Unix(clockOffTimeUnix, 0)
		fmt.Printf("%v\t\t%v\t%v\t%v\n", timeshiftID, projectID, clockOnTime, clockOffTime)
	}
	fmt.Println("==== projects ====")
	fmt.Println("project_id\tnamespace_id\tname")
	projects, err := db.Query("SELECT * FROM projects")
	if err != nil {
		panic(err)
	}
	for projects.Next() {
		var projectID int64
		var namespaceID sql.NullInt64
		var name string
		projects.Scan(&projectID, &namespaceID, &name)
		fmt.Printf("%v\t\t%v\t%v\n", projectID, namespaceID, name)
	}
	fmt.Println("==== namespaces ====")
	fmt.Println("namespace_id\tname")
	namespaces, err := db.Query("SELECT * FROM namespaces")
	if err != nil {
		panic(err)
	}
	for namespaces.Next() {
		var namespaceID int64
		var name string
		namespaces.Scan(&namespaceID, &name)
		fmt.Printf("%v\t\t%v\n", namespaceID, name)
	}
	return nil
}

// helper function finds namespaceId from namespace
func getNamespaceID(db *sql.DB, namespace string) (exists bool, id int64) {
	queryString := `
    SELECT namespace_id from namespaces
		WHERE name=? 
		LIMIT 1`
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
			timeshift_id INTEGER PRIMARY KEY,
			project_id INTEGER,
			clock_on_time INTEGER,
			clock_off_time INTEGER,
			FOREIGN KEY (project_id) REFERENCES projects(project_id)
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

// helper function formats the namespace + name output of a timeshift
func formatName(shift timeshift) string {
	if shift.project.namespace == "" {
		return shift.project.name
	} else {
		return shift.project.namespace + "." + shift.project.name
	}
}

type project struct {
	name      string
	namespace string
}

type timeshiftsDAO struct {
	dbDriver   string
	dbFilepath string
}

type timeshift struct {
	project      project
	clockOnTime  time.Time
	clockOffTime time.Time
}

type timeshiftQuery struct {
	projectName string
	namespace   string
	from        time.Time
	to          time.Time
}

type ErrTimeshiftAlreadyRunning timeshift

func (e ErrTimeshiftAlreadyRunning) Error() string {
	originalShift := timeshift(e)
	outputStr := "There is already a running timeshift for project %v: previous timeshift started at %v.\n"
	return fmt.Sprintf(outputStr, formatName(originalShift), originalShift.clockOnTime)
}
