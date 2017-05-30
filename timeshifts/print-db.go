package timeshifts

import (
	"database/sql"
	"fmt"
	"time"
)

// DEBUG
func (dbRef DB) PrintDB() error {
	db, err := sql.Open(dbRef.DbDriver, dbRef.DbFilepath)
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
