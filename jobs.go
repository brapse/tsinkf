package main

import (
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

func (job Job) Content() string {
  // XXX refactor this
  return ReadFile(store.getPath(NEW,job.hash))
}

func (job Job) LastTouch() time.Time {
  return LastTouch(store.getPath(NEW, job.hash))
}

func (job *Job) ToString() string {
	return strings.Join([]string{
          FormatTime(job.LastTouch()),
					STATELABELS[job.state],
					job.cmd,
					job.hash}, "\t")
}

func (job *Job) Run() error {
  // XXX: This sucks
  // The problem here is that it couple the execution with
  // responsibilities of the "store".
  // One way around this might be to use some named pipe or something
  // but there is a deeper problem of "running" and "storing" being coupled.

  // XXX: it might make sense to put the actual file in a stateless /jobs/ file
  outputFile := *baseDir + "/new/" + job.hash
  // XXX: This should append, so reruns use the sae file
  execution  := job.cmd + " &>" + outputFile

  //TODO setup some proper piping to clean up the process tree 
  cmd := exec.Command("bash","-c", execution)
  err := cmd.Run()

  return err
}

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
