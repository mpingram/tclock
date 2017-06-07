package main

import (
	"errors"
)

var ErrEmptyProject = errors.New("tclock: Project is empty or uninitialized.")

const shortTimeFormat = "Jan 2 3:04pm"
