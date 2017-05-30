package main

import (
	"fmt"
	"github.com/mpingram/tclock/timeshifts"
	"gopkg.in/urfave/cli.v1"
	"os"
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

	log := logger{}
	shortTimeFormat := "Jan 2 3:04pm"

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
				proj, err := parseProject(c.Args().First())
				if err != nil {
					return err
				}
				shift, err := timeshiftsDB.ClockOn(proj, forceOverwrite)
				if err != nil {
					log.e(err.Error())
					return err
				}
				log.i("Project %v clocked on: %v", c.Args().First(), shift.ClockOnTime.Format(shortTimeFormat))
				return nil
			},
		},
		{
			Name:  "off",
			Usage: "End a timeshift for the specified project.",
			Action: func(c *cli.Context) error {
				noProject := false
				proj, err := parseProject(c.Args().First())
				switch {
				case err == ErrEmptyProject:
					noProject = true
				case err != nil:
					return err
				}
				// get a list of all running timeshifts
				runningShifts, err := timeshiftsDB.GetRunningShifts()
				switch {
				// if there are no running timeshifts, throw err
				case err == timeshifts.ErrNoTimeshifts:
					log.e(err.Error())
					return nil
				case err != nil:
					log.e(err.Error())
					return err
				// if there is one running timeshift, and no specific project
				//  was passed to tclock off, clock off that project
				case len(runningShifts) == 1 && noProject == true:
					shift := runningShifts[0]
					shift, err = timeshiftsDB.ClockOff(shift.Project)
					if err != nil {
						return err
					}
					log.i("Project %v clocked off: %v, ran for %v",
						timeshifts.FormatProject(shift.Project),
						shift.ClockOffTime.Format(shortTimeFormat),
						shift.ClockOffTime.Sub(shift.ClockOnTime),
					)
					return nil
					// if there are multiple running timeshifts and no
					//   specific project was passed to tclock off, print a
					//   list of all running timeshifts
				case len(runningShifts) > 1 && noProject == true:
					outputStr := "Multiple running timeshifts:\n"
					outputStr += listTimeshifts(runningShifts)
					outputStr += "Clock off one of these shifts by calling tclock off <project>"
					log.i(outputStr)
					return nil
				// if a project has been passed to tclock off, clock off that project
				//   if it's running
				default:
					shift, err := timeshiftsDB.ClockOff(proj)
					switch {
					case err == timeshifts.ErrNoTimeshifts:
						outputStr := fmt.Sprintf("No running timeshift for project %v.\n", timeshifts.FormatProject(proj))
						outputStr += "Here is a list of all running timeshifts:\n"
						outputStr += listTimeshifts(runningShifts)
						log.i(outputStr)
						return nil
					case err != nil:
						return err
					default:
						log.i("Shift %v clocked out: %v. Ran for %v [%v - %v]",
							timeshifts.FormatProject(proj),
							shift.ClockOffTime.Format(shortTimeFormat),
							shift.ClockOnTime.Format("03:04pm"),
							shift.ClockOffTime.Format("03:04pm"),
						)
					}
					return nil
				}
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
