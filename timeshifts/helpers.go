package timeshifts

import (
	"database/sql"
	"time"
)

//  formats the namespace + name output of a timeshift
func FormatProject(project Project) string {
	if project.Namespace == "" {
		return project.Name
	} else {
		return project.Namespace + "." + project.Name
	}
}

//  checks if there is an unclosed timeshift for shift's project
func isTimeshiftAlreadyRunning(db *sql.DB, project Project) (running bool, prevShift Timeshift, err error) {
	exists, projectID, err := getProjectID(db, project)
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
			return false, prevShift, nil
		case err != nil:
			return false, prevShift, err
		default:
			var proj Project
			if maybeNamespace.Valid == true {
				namespace := maybeNamespace.String
				proj = Project{Name: name, Namespace: namespace}
			} else {
				proj = Project{Name: name}
			}
			clockOnTime := time.Unix(clockOnTimeUnix, 0)
			prevShift = Timeshift{Project: proj, ClockOnTime: clockOnTime}
			return true, prevShift, nil
		}
	}
}

//  inserts new namespace into db
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

//  inserts new project into db
func addProject(db *sql.DB, project Project) (int64, error) {
	createProjectStmt, err := db.Prepare(`INSERT INTO projects(name, namespace_id, project_id) VALUES (?,?,NULL)`)
	var (
		res sql.Result
	)
	if err != nil {
		return 0, err
	}
	projectHasNamespace := project.Namespace != ""
	if projectHasNamespace == true {
		namespaceExists, namespaceID := getNamespaceID(db, project.Namespace)
		if namespaceExists == false {
			// reassign namespaceID to newly created value
			namespaceID, err = addNamespace(db, project.Namespace)
			if err != nil {
				return 0, err
			}
		}
		res, err = createProjectStmt.Exec(project.Name, namespaceID)
	} else {
		res, err = createProjectStmt.Exec(project.Name, nil)
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

//  finds projectID from name and namespace
func getProjectID(db *sql.DB, proj Project) (exists bool, id int64, err error) {
	hasNamespace := proj.Namespace != ""
	if hasNamespace == false {
		queryString := `
			SELECT project_id FROM projects
			WHERE projects.name=? AND projects.namespace_id IS NULL`
		err = db.QueryRow(queryString, proj.Name).Scan(&id)
	} else {
		queryString := `
				SELECT project_id FROM projects 
					INNER JOIN namespaces ON projects.namespace_id = namespaces.namespace_id 
				WHERE projects.name=? AND namespaces.name=?`
		err = db.QueryRow(queryString, proj.Name, proj.Namespace).Scan(&id)
	}
	switch {
	case err == sql.ErrNoRows:
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

//  finds namespaceId from namespace
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
