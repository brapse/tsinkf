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
	force := c.Flags.Bool("force", false, "operate on all jobs")

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

	for _, job := range jobList {
		if len(jobIDs) == 0 || contains(job.id, jobIDs) > -1 {
			if !*force {
				fmt.Println(job.ToString())
				fmt.Printf("Reset to %s, \"yes\" to confirm: ", STATELABELS[NEWSTATE])
				var input string
				_, err := fmt.Scanf("%s", &input)
				if err != nil {
					panic(err)
				}

				if input == "yes" {
					job.SetState(NEWSTATE)
				}
			} else {
				job.SetState(NEWSTATE)
			}
		}
	}

	return 0
}

func init() {
	cmd := NewCmd("reset", "show the status of commands", "tsinkf reset [-v] [taskIDs]", resetFn)
	cmdList[cmd.Name] = cmd
}
