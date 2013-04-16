package main
import (
  "os"
)

// Storage stuff
type Store struct {
  baseDir string
}

func NewStore(relPath string) *Store {
  // TODO: check if we need to add the full path
  cwd, err := os.Getwd()
  if err != nil {
    panic(err)
  }

  baseDir := cwd + "/" + relPath

  err = os.MkdirAll(baseDir, 0755)
  if err != nil {
    panic(err)
  }

  // check and establish a lock
  lockFile := baseDir + "/lock"
  if FileExists(lockFile) {
    panic(baseDir + "is locked!\n can't run locked directory")
  } else {
    TouchFile(lockFile)
  }

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
  if (jobState == UNKNOWN) {
    panic("Don't store UNKNOWN state")
  }

  return s.baseDir + "/" + STATELABELS[jobState] + "/" + jobHash
}

func (s Store) get(key string) JobState {
  if !FileExists(s.getPath(NEW, key)) {
    return UNKNOWN
  }

  if FileExists(s.getPath(RUNNING, key)) {
    return RUNNING
  } else if FileExists(s.getPath(FAILED, key)) {
    return FAILED
  } else if FileExists(s.getPath(SUCCEEDED,key)) {
    return SUCCEEDED
  }

  return NEW
}

func (s *Store) set(jobHash string, jobState JobState) {
  TouchFile(s.getPath(NEW, jobHash))
  // Delete previous state files
  for _, state := range []JobState{RUNNING,SUCCEEDED,FAILED} {
    RemoveFile(s.getPath(state, jobHash))
  }

  if jobState != NEW {
    err := os.Symlink(s.getPath(NEW, jobHash), s.getPath(jobState, jobHash))
    if err != nil {
      panic(err)
    }
  }

}

func (s *Store) close() {
  RemoveFile(s.baseDir + "/lock")
}
