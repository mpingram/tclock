package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "tclock"
	app.Usage = "Record the time you spend working on projects"
	app.Commands = []cli.Command{
		{
			Name:  "on",
			Usage: "Start a timeshift for the specified project.",
			Action: func(c *cli.Context) error {
				clockOnTime := time.Now()
				projectName, namespace := parseProject(c.Args().First())

				project := projectStruct{projectName, namespace}
				clockOnDataMessage := timeslotDataMessage{project: project, clockOnTime: clockOnTime}
				sendData(clockOnDataMessage)
				return nil
			},
		},
		{
			Name:  "off",
			Usage: "End a timeshift for the specified project.",
			Action: func(c *cli.Context) error {
				clockOffTime := time.Now()
				projectName, namespace := parseProject(c.Args().First())
				project := projectStruct{projectName, namespace}
				clockOffDataMessage := timeslotDataMessage{project: project, clockOffTime: clockOffTime}
				sendData(clockOffDataMessage)
				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Show the currently active timeshifts and how long they've been running.",
			Action: func(c *cli.Context) error {
				fmt.Println("Implement me.")
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

func parseProject(fullProjectStr string) (projectName, namespace string) {
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
	return
}
