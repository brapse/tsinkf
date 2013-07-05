package main

import (
  "fmt"
)

var versionFn = func(c *Cmd, args []string) int {
  fmt.Printf("tsinkf version %s\n", Version)
  return 0
}

func init() {
  cmd := NewCmd("version", "display version and exit", "tsinkf version", versionFn)
  cmdList[cmd.Name] = cmd
}
