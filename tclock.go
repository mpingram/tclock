package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
)

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

func main() {
	app := cli.NewApp()
	app.Name = "tclock"
	app.Usage = "Record the time you spend working on projects"
	app.Commands = []cli.Command{
		{
			Name:  "on",
			Usage: "Begin logging time for the specified project. If no project is passed, an unnamed project is begun",
			Action: func(c *cli.Context) error {

				clockOnTime := time.Now()
				projectName, namespace := parseProject(c.Args().First())
				fmt.Printf("Clocked on project %s at %s\n", projectName, clockOnTime)

				project := projectStruct{projectName, namespace}
				clockOnDataMessage := timeslotDataMessage{project: project, clockOnTime: clockOnTime}
				sendData(clockOnDataMessage)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
