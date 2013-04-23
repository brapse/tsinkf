package main
import (
  "os"
  "time"
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

  lockFile := baseDir + "/lock"
  if FileExists(lockFile) {
    panic(baseDir + "is locked!\n can't run locked directory")
  } else {
    TouchFile(lockFile)
  }

  for _, directory := range []string{"NEW", "RUNNING", "FAILED", "SUCCEEDED"} {
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

func (s Store) GetState(jobID string) JobState {
  if !FileExists(s.getPath(NEW, jobID)) {
    return UNKNOWN
  }

  if FileExists(s.getPath(RUNNING, jobID)) {
    return RUNNING
  } else if FileExists(s.getPath(FAILED, jobID)) {
    return FAILED
  } else if FileExists(s.getPath(SUCCEEDED,jobID)) {
    return SUCCEEDED
  }

  return NEW
}

func (s Store) SetOutput(jobID string, output string) {
  AppendToFile(s.getPath(NEW,jobID), output)
}

func (s Store) GetOutput(jobID string) string {
  return ReadFile(s.getPath(NEW, jobID))
}

func (s Store) GetLastTouch(jobID string) time.Time {
  return LastTouch(s.getPath(NEW,jobID))
}

func (s Store) GetJobIDs() (result []string) {
	for _, filename := range ListFiles(s.baseDir + "/NEW") {
		result = append(result, filename)
	}

	return result
}

func (s Store) SetState(jobID string, jobState JobState) {
  TouchFile(s.getPath(NEW, jobID))
  // Delete previous state files
  for _, state := range []JobState{RUNNING,SUCCEEDED,FAILED} {
    RemoveFile(s.getPath(state, jobID))
  }

  if jobState != NEW {
    err := os.Symlink(s.getPath(NEW, jobID), s.getPath(jobState, jobID))
    if err != nil {
      panic(err)
    }
  }
}

func (s Store) Reset() {
  for _, jobID := range s.GetJobIDs() {
    for _, state := range []JobState{RUNNING,SUCCEEDED,FAILED} {
      RemoveFile(s.getPath(state, jobID))
    }
  }
}

func (s Store) Close() {
  RemoveFile(s.baseDir + "/lock")
}
