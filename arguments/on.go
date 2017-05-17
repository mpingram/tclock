package arguments

import (
  "github.github.com/jessevdk/go-flags"
  "fmt"
)

type ClockOnCommand struct {
  ClockOn bool `long:"on" description "clock on; begin logging time on specified project. If no project is passed, an unnamed project will be started."`
}

var clockOnCommand ClockOnCommand

func (x *ClockOnCommand) Execute(args []string) error {
  fmt.Printf("Clocked on.")
  return nil
}

func init() {
  parser.AddCommand("on",
    "Clock on",
    "Begin logging time on specified project. If no project is passed, and unnamed project will be started."
    &ClockOnCommand)
}
