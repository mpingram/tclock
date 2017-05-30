package main

import (
	"github.com/mpingram/tclock/timeshifts"
	"strings"
)

func parseProject(fullProjectStr string) (timeshifts.Project, error) {
	var projectName, namespace string
	splitName := strings.SplitN(fullProjectStr, ".", 2)
	if len(splitName) > 1 {
		namespace = splitName[0]
		projectName = splitName[1]
		return timeshifts.Project{projectName, namespace}, nil
	} else if len(splitName) == 1 {
		namespace = ""
		projectName = fullProjectStr
		return timeshifts.Project{projectName, ""}, nil
	} else {
		return timeshifts.Project{}, ErrEmptyProject
	}
}
