package main

import (
	"fmt"
	"github.com/mpingram/tclock/timeshifts"
	"os"
	"text/tabwriter"
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

func printTimeshifts(shifts []timeshifts.Timeshift) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	for i, shift := range shifts {
		line := fmt.Sprintf("  [%v]\t%v\t started %v \t [%v]",
			i,
			timeshifts.FormatProject(shift.Project),
			shift.ClockOnTime.Format(shortTimeFormat),
			durFormat(time.Since(shift.ClockOnTime)),
		)
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
