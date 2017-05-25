package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strings"
	"time"
)

func main() {

	// TODO: read this in from config
	db_type := "sqlite3"
	db_location := "./timeshifts.db"
	timeshiftsDB := timeshiftsDAO{db_type, db_location}

	err := timeshiftsDB.init()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "tclock"
	app.Usage = "Record the time you spend working on projects"
	app.Commands = []cli.Command{
		{
			Name:  "test",
			Usage: "test",
			Action: func(c *cli.Context) error {
				err := timeshiftsDB.init()
				fmt.Println(err)
				return nil
			},
		},

		{
			Name:  "on",
			Usage: "Start a timeshift for the specified project.",
			Action: func(c *cli.Context) error {
				clockInTime := time.Now()
				proj := parseProject(c.Args().First())
				shift := timeshift{project: proj, clockInTime: clockInTime}
				err := timeshiftsDB.clockIn(shift)
				if err != nil {
					printErr(err)
					// DEV ONLY
					panic(err)
					//return err
				}
				return nil
			},
		},
		{
			Name:  "off",
			Usage: "End a timeshift for the specified project.",
			Action: func(c *cli.Context) error {
				clockOutTime := time.Now()
				proj := parseProject(c.Args().First())
				timeshiftClockOut := timeshift{project: proj, clockOutTime: clockOutTime}
				err := timeshiftsDB.clockOut(timeshiftClockOut)
				if err != nil {
					printErr(err)
				}
				return err
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Show the currently active timeshifts and how long they've been running.",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "report",
			Usage: "Show a report of all timeshifts from the past week [default]. To view timeshifts for a certain project or namespace, pass that project or namespace as the first argument. To view timeshifts over a different period, pass the [-d]--duration flag set to the duration to include, which should be a number immediately followed by a unit of time, one of [(d)ay, (w)eek, (m)onth, (y)ear]. \nFor exmample, to display the past two weeks of timeshifts spent in namespace world-domination, type \n\ttclock report world-domination -d=2w\n",
			Action: func(c *cli.Context) error {
				fmt.Println("Implement me.")
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func parseProject(fullProjectStr string) project {
	var projectName, namespace string
	splitName := strings.SplitN(fullProjectStr, ".", 2)
	if len(splitName) > 1 {
		namespace = splitName[0]
		projectName = splitName[1]
	} else if len(splitName) == 1 {
		namespace = ""
		projectName = fullProjectStr
	} else {
		namespace = ""
		projectName = "unnamed"
	}
	return project{projectName, namespace}
}

func printErr(err error) {
	fmt.Println(err)
}
