package main

import (
  "os/exec"
	"strings"
)

type JobState int

const (
  UNKNOWN   JobState = -1
  NEW       JobState = 0
  RUNNING   JobState = 1
  FAILED    JobState = 2
  SUCCEEDED JobState = 3
)

var STATELABELS = map[JobState] string {
  UNKNOWN: "UNKNOWN",
  NEW: "NEW",
  RUNNING: "RUNNING",
  FAILED: "FAILED",
  SUCCEEDED: "SUCCEEDED",
}

type Job struct {
  hash string
  cmd string
  state JobState
}

func NewJob(cmd string) *Job {
  hash := CreateHash(cmd)

  job := &Job{hash, cmd, UNKNOWN}
  return job
}

func (job *Job) ToString() string {
	return strings.Join([]string{
		// time label
					STATELABELS[job.state],
					job.cmd,
					job.hash}, "\t")
}

func (job *Job) Run() error {
  // XXX: This sucks
  outputFile := *baseDir + "/new/" + job.hash
  execution  := job.cmd + " &>" + outputFile

  //TODO setup some proper piping to clean up the process tree 
  cmd := exec.Command("bash","-c", execution)
  err := cmd.Run()

  return err
}

//jobList stuff

type JobList []Job

var journal = &Journal{}
var store = &Store{}

func NewJobList(stor *Store, jrnl *Journal) JobList {
  store = stor
  journal = jrnl

  jobList := JobList{}

  for hash, state:= range store.GetAll() {
		job := NewJob(DecodeHash(hash))
		if state == RUNNING {
			state = FAILED
		}
    job.state = state

		jobList = append(jobList, *job)
  }

  return jobList
}

func (jobList JobList) Include(job *Job) bool {
	for _, j := range jobList {
		if job.hash == j.hash {
			return true
		}
	}
	return false
}

// TODO: store state in Job and not JobList
func (jobList *JobList) Add(job Job) {
		store.Set(job.hash, job.state)
    foo := append(*jobList, job)
    *jobList = foo
}

func (jobList *JobList) Update(job *Job, newState JobState) {
  journal.Log(*job, newState)
  store.Set(job.hash, newState)
  job.state = newState
}
