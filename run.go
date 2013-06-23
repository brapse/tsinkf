package main

import (
  "strings"
)

var runFn CmdFn  = func(c * Cmd, args []string) int {
  baseDir := c.Flags.String("d", ".tsinkf", "base directory")
  verbose := c.Flags.Bool("v", false, "verbose output")

  c.Flags.Parse(args)

  if len(args) == 0 {
    c.Usage()
    return 1
  }

  leftover := c.Flags.Args()

	store := NewStore(*baseDir)
	journal := NewJournal(*verbose, *baseDir+"/journal.log")
	defer store.Close()
	defer journal.Close()

	jobList := NewJobList(store, journal)

	cmd := strings.Join(leftover, " ")
	job := NewJob(cmd, *store, *journal)
	if !jobList.Include(*job) {
		job.SetState(NEW)
		jobList.Add(*job)
	}

	if job.GetState() == NEW {
		job.SetState(RUNNING)
		err := job.Run()
		if err == nil {
			job.SetState(SUCCEEDED)
		} else {
			job.SetState(FAILED)
		}
	}
  return 0
}


func init() {
  cmd := NewCmd("run", "run a command", "tsinkf run [-v] cmd...", runFn)
  cmdList[cmd.Name] = cmd
}
