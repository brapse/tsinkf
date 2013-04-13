package main

import (
  "os"
  "os/exec"
  "strings"
  "log"
  "flag"
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

// Journal

func init() {
  flag.Parse()
}

var debug = flag.Bool("v", false, "Debug info")

func Journal(job Job, toState JobState) {
    fromLabel := STATELABELS[job.state]
    toLabel   := STATELABELS[toState]
    msg := strings.Join([]string{
      fromLabel + "->" + toLabel,
      job.cmd }, "\t")
    if *debug {
        log.Printf(msg)
    }
    // TODO: Write to file
}

//filesystem stuff
func touchFile (filename string) {
  _, err := os.OpenFile(filename, os.O_CREATE, 0666)
  if err != nil {
    panic(err)
  }
}

func fileExists(filepath string) bool {
  _, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
  if err != nil {
    if os.IsNotExist(err) {
      return false
    } else {
      panic(err)
    }
  }

  return true
}

// Storage stuff
type Store struct {
  baseDir string
}

func newStore(relPath string) *Store {
  // TODO: check if we need to add the full path
  cwd, err := os.Getwd()
  if err != nil {
    panic(err)
  }

  baseDir := cwd + "/" + relPath

  for _, directory := range []string{"new", "running", "failed", "succeeded"} {
    todo := baseDir + "/" + directory
    err := os.MkdirAll(todo, 0755)
    if err != nil {
      panic(err)
    }
  }

  return &Store{baseDir}
}

func (s Store) getPath(jobState JobState, jobHash string) string {
  if jobState == NEW {
    return s.baseDir + "/new/" + jobHash
  } else if jobState == RUNNING {
    return s.baseDir + "/running/" + jobHash
  } else if jobState == FAILED {
    return s.baseDir + "/failed/" + jobHash
  }

  // XXX: this may prove to be wrong
  return s.baseDir + "/succeeded/" + jobHash
}

func (s Store) get(key string) JobState {
  if !fileExists(s.getPath(NEW, key)) {
    return UNKNOWN
  }

  if fileExists(s.getPath(RUNNING, key)) {
    return RUNNING
  } else if fileExists(s.getPath(FAILED, key)) {
    return FAILED
  } else if fileExists(s.getPath(SUCCEEDED,key)) {
    return SUCCEEDED
  }

  return NEW
}

func (s *Store) set(j Job, to JobState) {
  if j.state == UNKNOWN && to == NEW{
    touchFile(s.getPath(NEW, j.hash))
  } else if j.state == NEW && to == RUNNING {
    err := os.Symlink(s.getPath(NEW, j.hash), s.getPath(RUNNING, j.hash))
    if err != nil {
      panic(err)
    }
  } else if j.state == RUNNING && to == SUCCEEDED {
    err := os.Remove(s.getPath(RUNNING, j.hash))
    if err != nil {
      panic(err)
    }

    err = os.Symlink(s.getPath(NEW, j.hash), s.getPath(SUCCEEDED, j.hash))
    if err != nil {
      panic(err)
    }
  } else if j.state == RUNNING && to == FAILED {
    err := os.Remove(s.getPath(RUNNING, j.hash))
    if err != nil {
      panic(err)
    }

    err = os.Symlink(s.getPath(NEW, j.hash), s.getPath(FAILED, j.hash))
    if err != nil {
      panic(err)
    }
  }
}

//jobs stuff

type Job struct {
  hash string
  cmd string
  state JobState
  store *Store
}

func NewJob(cmd string, store *Store) *Job {
  hash := CreateHash(cmd)

  state := store.get(hash)

  job := &Job{hash, cmd, state, store}

  if state == UNKNOWN {
    job.update(NEW)
  }
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

func (job *Job) update(to JobState) {
  // XXX: Exceptions?
  Journal(*job, to)
  job.store.set(*job, to)
  job.state = to
}

//jobList stuff

type JobList []*Job

func NewJobList(cmdName string, listing []string, baseDir string) JobList {
  jobList := JobList{}
  store := newStore(baseDir)
  for _, args := range listing {
    cmd := cmdName + " " + args

    jobList = append(jobList,NewJob(cmd, store))
  }

  return jobList
}

func (jobList JobList) available() (availableJobs []*Job) {
  for _,job := range jobList {
    if job.state == NEW {
      availableJobs = append(availableJobs, job)
    }
  }

  return
}
