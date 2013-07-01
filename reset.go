package main

var resetFn CmdFn = func(c *Cmd, args []string) int {
	verbose := c.Flags.Bool("v", false, "verbose output")
	state := c.Flags.String("s", STATELABELS[NEW], "state to set")

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
	resetable := JobSpecific(jobIDs)

	for _, job := range jobList {
		if resetable(job) {
			job.SetState(NEWSTATE)
		}
	}
	return 0
}

func init() {
	cmd := NewCmd("reset", "show the status of commands", "tsinkf reset [-v] [taskIDs]", resetFn)
	cmdList[cmd.Name] = cmd
}
