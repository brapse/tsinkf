package main

import (
  "os/exec"
	"strings"
  "time"
  "io"
  "io/ioutil"
  "os"
  "sync"
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

var JOBSTATEIDS = map[string] JobState {
  STATELABELS[UNKNOWN]:   UNKNOWN,
  STATELABELS[NEW]:       NEW,
  STATELABELS[RUNNING]:   RUNNING,
  STATELABELS[FAILED]:    FAILED,
  STATELABELS[SUCCEEDED]: SUCCEEDED,
}

type Job struct {
  id string
  cmd string
  journal Journal
  store Store
}

func NewJob(cmd string, store Store, journal Journal) *Job {
  jobID := CreateHash(cmd)
  state := store.GetState(jobID)
  job := &Job{jobID, cmd, journal, store}

  if state == UNKNOWN {
    store.Setup(jobID, cmd)
    journal.Log(*job, state)
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

func (job *Job) Run() error {
  cmd := exec.Command("bash", "-c", job.cmd)
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    panic(err)
  }
  stderr, err := cmd.StderrPipe()
  if err != nil {
    panic(err)
  }

  tStdout := io.TeeReader(stdout, os.Stdout)
  tStderr := io.TeeReader(stderr, os.Stderr)

  reader, writer := io.Pipe()

  var wg sync.WaitGroup

  go func() {
    wg.Add(1)
    defer wg.Done()

    buf, err := ioutil.ReadAll(reader)
    if err != nil {
      panic(err)
    }
    job.store.SetOutput(job.id, string(buf))
  }()

  go func () {
    wg.Add(1)
    defer wg.Done()
    _, err := io.Copy(writer, tStdout)
    if err != nil {
      panic(err)
    }
  }()

  go func () {
    wg.Add(1)
    defer wg.Done()
    _, err := io.Copy(writer, tStderr)
    if err != nil {
      panic(err)
    }
  }()

  result := cmd.Run()
  /* XXX
  There could be a race condition here,
  as there is no guarentee that the child processes are done
  done writing to the pipe
  */
  writer.Close()
  wg.Wait()

  return result
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
		if job.id== j.id {
			return true
		}
	}
	return false
}

func (jobList *JobList) Add(job Job) {
    foo := append(*jobList, job)
    *jobList = foo
}
