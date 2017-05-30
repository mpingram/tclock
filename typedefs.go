package main

import (
	"errors"
)

var ErrEmptyProject = errors.New("tclock: Project is empty or uninitialized.")
