package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type JobState int

const (
	UNKNOWN   JobState = -1
	NEW       JobState = 0
	RUNNING   JobState = 1
	FAILED    JobState = 2
	SUCCEEDED JobState = 3
)

var STATELABELS = map[JobState]string{
	UNKNOWN:   "UNKNOWN",
	NEW:       "NEW",
	RUNNING:   "RUNNING",
	FAILED:    "FAILED",
	SUCCEEDED: "SUCCEEDED",
}

var JOBSTATEIDS = map[string]JobState{
	STATELABELS[UNKNOWN]:   UNKNOWN,
	STATELABELS[NEW]:       NEW,
	STATELABELS[RUNNING]:   RUNNING,
	STATELABELS[FAILED]:    FAILED,
	STATELABELS[SUCCEEDED]: SUCCEEDED,
}

type Job struct {
	id      string
	cmd     string
	journal Journal
	store   Store
}

func NewJob(cmd string, store Store, journal Journal) *Job {
	jobID := CreateHash(cmd)
	state := store.GetState(jobID)
	job := &Job{jobID, cmd, journal, store}

	if state == UNKNOWN {
		store.Setup(jobID, cmd, UNKNOWN)
	}

	return job
}

func (job Job) SetState(state JobState) {
	job.journal.Log(job, state)
	job.store.SetState(job.id, state)
}

func (job Job) GetOutput() string {
	return job.store.GetOutput(job.id)
}

func (job Job) GetLastTouch() time.Time {
	return job.store.GetLastTouch(job.id)
}

func (job Job) GetState() JobState {
	return job.store.GetState(job.id)
}

func (job *Job) ToString() string {
	return strings.Join([]string{
		FormatTime(job.GetLastTouch()),
		STATELABELS[job.GetState()],
		job.cmd,
		job.id}, "\t")
}

func (job *Job) Run() int {
	buf := bytes.Buffer{}
	defer func() { job.store.SetOutput(job.id, buf.String()) }()

	cmd := exec.Command("bash", "-c", job.cmd)
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)

	if err := cmd.Start(); err != nil {
    panic(err)
	}

	if err := cmd.Wait(); err != nil {
    return 1
	}

	return 0
}

type JobList []Job

var journal = &Journal{}
var store = &Store{}

func NewJobList(stor *Store, jrnl *Journal) JobList {
	store = stor
	journal = jrnl

	jobList := JobList{}

	for _, jobID := range store.GetJobIDs() {
		cmd := store.GetCmd(jobID)
		job := NewJob(cmd, *stor, *jrnl)
		if job.GetState() == RUNNING {
			job.SetState(FAILED)
		}

		jobList = append(jobList, *job)
	}

	return jobList
}

func (jobList JobList) Include(job Job) bool {
	for _, j := range jobList {
		if job.id == j.id {
			return true
		}
	}
	return false
}

func (jobList *JobList) Add(job Job) {
	foo := append(*jobList, job)
	*jobList = foo
}
