package main

import (
	"strings"
)

var runFn CmdFn = func(c *Cmd, args []string) int {
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

	store.Lock()
	defer store.Close()
	defer journal.Close()

	cmd := strings.Join(leftover, " ")
	job := NewJob(cmd, *store, *journal)

	runJob := func() int {
		job.SetState(RUNNING)
		status := job.Run()
		if status == 0 {
			job.SetState(SUCCEEDED)
		} else {
			job.SetState(FAILED)
		}

		return status
	}

	switch job.GetState() {
	case UNKNOWN:
		job.SetState(NEW)
		return runJob()
	case NEW:
		return runJob()
	case RUNNING:
		job.SetState(FAILED)
		return 1
	case FAILED:
		return 1
	case SUCCEEDED:
		return 0
	}
	return 1

}

func init() {
	cmd := NewCmd("run", "run a command", "tsinkf run [-v] cmd...", runFn)
	cmdList[cmd.Name] = cmd
}
