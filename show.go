package main

import (
  "fmt"
)

var showFn CmdFn  = func(c * Cmd, args []string) int {
  verbose := c.Flags.Bool("v", false, "verbose output")

  c.Flags.Parse(args)

	store := NewStore(*root)
	journal := NewJournal(*verbose, *root+"/journal.log")

	defer store.Close()
	defer journal.Close()

	jobList := NewJobList(store, journal)

  jobIDs := c.Flags.Args()

	printable := JobSpecific(jobIDs)

	for _, job := range jobList {
		if printable(job) {
			if ! *verbose {
				fmt.Println(job.ToString())
			} else {
				fmt.Println(job.ToString() + "\n" + job.GetOutput())
			}
		}
	}

  return 0
}

func init() {
  cmd := NewCmd("show", "show the status of commands", "tsinkf show[-v] [taskID]", showFn)
  cmdList[cmd.Name] = cmd
}
