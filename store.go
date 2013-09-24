package main

import (
	"os"
	"time"
)

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

	jobsDir := baseDir + "/jobs"
	err = os.MkdirAll(jobsDir, 0755)
	if err != nil {
		panic(err)
	}

	return &Store{baseDir}
}

func (s Store) Lock() {
	lockFile := s.baseDir + "/lock"
	if FileExists(lockFile) {
		panic(s.baseDir + "is locked!\n can't run locked directory")
	} else {
		TouchFile(lockFile)
	}
}

func (s Store) Unlock() {
	lockFile := s.baseDir + "/lock"
	if !FileExists(lockFile) {
		panic(s.baseDir + "is not locked!\n can't unlock")
	} else {
		DeleteFile(lockFile)
	}
}

func (s Store) getPath(jobID string, field string) string {
	return s.baseDir + "/jobs/" + jobID + "/" + field
}

func (s Store) Setup(jobID string, cmd string, state JobState) {
	jobsDir := s.baseDir + "/jobs/" + jobID
	err := os.MkdirAll(jobsDir, 0755)
	if err != nil {
		panic(err)
	}

	SetFile(s.getPath(jobID, "cmd"), cmd)
	TouchFile(s.getPath(jobID, "output"))
	s.SetState(jobID, state)
}
func (s Store) GetState(jobID string) JobState {
	stateFile := s.getPath(jobID, "state")

	if !FileExists(stateFile) {
		return UNKNOWN
	}
	contents := ReadFile(stateFile)
	return JOBSTATEIDS[contents]
}

func (s Store) SetOutput(jobID string, output string) {
	AppendToFile(s.getPath(jobID, "output"), output)
}

func (s Store) GetOutput(jobID string) string {
	return ReadFile(s.getPath(jobID, "output"))
}

func (s Store) GetLastTouch(jobID string) time.Time {
	return LastTouch(s.getPath(jobID, "state"))
}

func (s Store) GetCmd(jobID string) string {
	return ReadFile(s.getPath(jobID, "cmd"))
}

func (s Store) GetJobIDs() (result []string) {
	for _, filename := range ListFiles(s.baseDir + "/jobs") {
		result = append(result, filename)
	}

	return result
}

func (s Store) SetState(jobID string, jobState JobState) {
	SetFile(s.getPath(jobID, "state"), STATELABELS[jobState])
}

func (s Store) Reset() {
	for _, jobID := range s.GetJobIDs() {
		s.SetState(jobID, NEW)
	}
}

func (s Store) Close() {
	RemoveFile(s.baseDir + "/lock")
}
