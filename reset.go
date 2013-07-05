package main

import (
	"fmt"
)

func contains(needle string, heystack []string) int {
	for i, possible := range heystack {
		if needle == possible {
			return i
		}
	}

	return -1
}

var resetFn CmdFn = func(c *Cmd, args []string) int {
	verbose := c.Flags.Bool("v", false, "verbose output")
	state := c.Flags.String("state", STATELABELS[NEW], "state to set")
	all := c.Flags.Bool("all", false, "operate on all jobs")

	c.Flags.Parse(args)

	NEWSTATE, found := JOBSTATEIDS[*state]
	if !found {
		panic("State not found:" + *state)
	}

	store := NewStore(*root)
	journal := NewJournal(*verbose, *root+"/journal.log")

	defer journal.Close()
	defer store.Close()

	jobIDs := c.Flags.Args()

	jobList := NewJobList(store, journal)

	if *all {

		for _, job := range jobList {
			fmt.Println(job.ToString())
		}

		fmt.Printf("Resetting all jobs to %s, type \"yes\" to confirm: ", STATELABELS[NEWSTATE])
		var input string
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			panic(err)
		}

		if input != "yes" {
			fmt.Println("ABORT")
			return 1
		}
	}

	for _, job := range jobList {
		if *all || contains(job.id, jobIDs) > -1 {
			job.SetState(NEWSTATE)
		}
	}

	return 0
}

func init() {
	cmd := NewCmd("reset", "show the status of commands", "tsinkf reset [-v] [taskIDs]", resetFn)
	cmdList[cmd.Name] = cmd
}
