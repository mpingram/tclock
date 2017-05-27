package main

import (
	"fmt"
	"github.com/mpingram/tclock/timeshifts"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strings"
	"time"
)

func main() {

	// TODO: read this in from config
	// TODO: custom error wrapper
	db_type := "sqlite3"
	db_location := "./timeshifts.db"
	timeshiftsDB := timeshifts.DB{db_type, db_location}

	err := timeshiftsDB.Init()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "tclock"
	app.Usage = "Record the time you spend working on projects"
	app.Commands = []cli.Command{
		{
			Name:    "report",
			Aliases: []string{"r"},
			Usage:   "Show active and previous timeshifts.",
			Action: func(c *cli.Context) error {
				err := timeshiftsDB.PrintDB()
				if err != nil {
					panic(err)
				}
				return nil
			},
		},

		{
			Name:  "on",
			Usage: "Start a timeshift for the specified project.",
			Action: func(c *cli.Context) error {
				forceOverwrite := false
				clockOnTime := time.Now()
				proj := parseProject(c.Args().First())
				shift := timeshifts.Timeshift{Project: proj, ClockOnTime: clockOnTime}
				err := timeshiftsDB.ClockOn(shift, forceOverwrite)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "off",
			Usage: "End a timeshift for the specified project.",
			Action: func(c *cli.Context) error {
				clockOffTime := time.Now()
				proj := parseProject(c.Args().First())
				timeshiftClockOff := timeshifts.Timeshift{Project: proj, ClockOffTime: clockOffTime}
				err := timeshiftsDB.ClockOff(timeshiftClockOff)
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

func parseProject(fullProjectStr string) timeshifts.Project {
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
	return timeshifts.Project{projectName, namespace}
}

func printErr(err error) {
	fmt.Println(err)
}
