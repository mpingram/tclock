package main

import (
	"fmt"
	"github.com/mpingram/tclock/timeshifts"
	"time"
)

func listTimeshifts(shifts []timeshifts.Timeshift) string {
	outputStr := ""
	for i, shift := range shifts {
		if shift.ClockOffTime.IsZero() {
			outputStr += fmt.Sprintf("%v: project %v, running for %v [%v - ]\n",
				i,
				timeshifts.FormatProject(shift.Project),
				time.Since(shift.ClockOnTime),
				shift.ClockOnTime.Format("03:04pm"),
			)
		} else {
			outputStr += fmt.Sprintf("%v: project %v, ran for %v [%v - %v]\n",
				i,
				timeshifts.FormatProject(shift.Project),
				shift.ClockOffTime.Sub(shift.ClockOnTime),
				shift.ClockOnTime.Format("03:04pm"),
				shift.ClockOffTime.Format("03:04pm"),
			)
		}
	}
	return outputStr
}
