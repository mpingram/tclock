package main

import (
  "os"
  "fmt"
  "./arguments"
)


func main() {
  args := os.Args

  fmt.Printf("Argument: %s\n", args[1])

  mainArg := args[1]
  
}
