package main

import (
  "os"
  "os/exec"
  "strings"
  "fmt"
  "hashing"
)

type JobState int

const (
  UNKNOWN   JobState = -1
  NEW       JobState = 0
  RUNNING   JobState = 1
  FAILED    JobState = 2
  SUCCEEDED JobState = 3
)

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

func newStore(baseDir string) *Store {
  // create directory structure that will be required for writting
  for _, directory := range []string{"new", "running", "failed", "succeeded"} {
    todo := *baseDir + "/" + directory
    err := os.MkdirAll(todo, 0755)
    if err != nil {
      panic(err)
    }
  }

  return Store{baseDir}
}

func (s Store) getPath(j Job, jobHash string) string {
  if jobState == NEW {
    return s.baseDir + "/new/" + jobHash
  } else if jobState == RUNNING {
    return s.baseDir + "/running/" + jobHash
  } else if s.jobState == FAILED {
    return s.baseDir + "/failed/" + jobHash
  }

  // XXX: this may prove to be wrong
  return s.baseDir + "/succeeded/" + j.hash
}

func (s Store) get(key string) jobState {
  if !fileExists(s.getPath(NEW, key)) {
    return UNKNOWN, fmt.Errorf("Job file does not exist")
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
  cwd, err := os.Getwd()

  if err != nil {
    panic(err)
  }

  if j.state == NEW && to == RUNNING {
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
  body string
  state JobState
  store Store
}

func NewJob(cmd string, store *Store) *Job {
  hash := hashing.CreateHash(cmd)

  // TODO: Setup storage
  state := store.get(hash)
  return &Job{hash, cmd, state, store}
}


func (job *Job) run() error {
  outputFile := *baseDir + "/new/" + job.hash
  execution  := job.body + " &>" + outputFile

  //TODO setup some proper piping to clean up the process tree 
  cmd := exec.Command("bash","-c", execution)
  err := cmd.Run()

  return err
}

func (job *Job) update(to JobState) {
  // XXX: Exceptions?
  store.set(job, to)
}

//jobList stuff

type JobList []*Job

func NewJobList(cmdName string, listing []string) JobList {
  jobList := JobList{}
  for _, args := range listing {
    cmd := cmdName + " " + args

    // SYNCING
    job := newJob(cmd, &store)
  }

  return jobList
}
