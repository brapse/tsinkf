package main

import (
  "os/exec"
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

//jobs stuff

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


func (job *Job) run() error {
  outputFile := *baseDir + "/new/" + job.hash
  execution  := job.cmd + " &>" + outputFile

  //TODO setup some proper piping to clean up the process tree 
  cmd := exec.Command("bash","-c", execution)
  err := cmd.Run()

  return err
}

//jobList stuff

type JobList []*Job

var journal = Journal{}
var store = Store{}

func NewJobList(cmdName string, listing []string, stor Store, jrnl Journal) JobList {
  store = stor
  journal = jrnl

  jobList := JobList{}

  for _, args := range listing {
    if len(args) > 0 {
      cmd := cmdName + " " + args
      job := NewJob(cmd)
      // sync state with store
      state := store.get(job.hash)
      if state == UNKNOWN {
        state = NEW
      } else if state == RUNNING {
        state = FAILED
      }

      store.set(job.hash, state)
      job.state = state

      jobList = append(jobList, job)
    }
  }

  return jobList
}

func (jobList JobList) update(job *Job, newState JobState) {
  // Update the journal
  journal.log(*job, newState)
  // Update the store
  store.set(job.hash, newState)
  job.state = newState
}

func (jobList JobList) available() (availableJobs []*Job) {
  for _,job := range jobList {
    if job.state == NEW {
      availableJobs = append(availableJobs, job)
    }
  }

  return
}
