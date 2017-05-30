package main

import (
	"fmt"
)

type logger struct{}

const basePrefix string = "[tclock]"

// info-level logging.
// Acts like fmt.printf(<string>, <value>, ... <value>)
// Appends \n to end of every string if not already there.
func (l logger) i(fmtString string, values ...interface{}) {
	prefix := basePrefix
	logPrint(prefix+fmtString, values...)
}

// error-level logging
// Acts like fmt.printf(<string>, <value>, ... <value>)
// Appends \n to end of every string if not already there.
func (l logger) e(fmtString string, values ...interface{}) {
	prefix := basePrefix + " ERROR: "
	logPrint(prefix+fmtString, values...)
}

// debug-level logging
// Acts like fmt.printf(<string>, <value>, ... <value>)
// Appends \n to end of every string if not already there.
func (l logger) d(fmtString string, values ...interface{}) {
	prefix := basePrefix + " DEBUG: "
	logPrint(prefix+fmtString, values...)
}

func logPrint(fmtString string, values ...interface{}) {
	lastChar := fmtString[len(fmtString)-1]
	if lastChar != '\n' {
		fmtString = fmtString + "\n"
	}
	fmt.Printf(fmtString, values...)
}
