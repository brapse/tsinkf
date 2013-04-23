package main

import (
  "flag"
  "os/exec"
  "strings"
  "bytes"
  "log"
  "fmt"
)

var baseDir = flag.String("dir", ".tsinkf", "directory where state files are created")

func init() {
  flag.Parse()
}

func getFrom(fromCmd string) (res []string) {
  cmd := exec.Command("bash", "-c", fromCmd)
  var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()
  if err != nil {
    log.Fatal(err)
  }
  for _, line := range strings.Split(out.String(), "\n") {
    if len(line) > 0 {
      res = append(res, line)
    }
  }
  return
}

// Run
type Run struct {
  from *string
  to *string
  verbose *bool
}

func (cmd *Run) DefineFlags(fs *flag.FlagSet) {
  cmd.from     = fs.String("from", "", "From command line tool")
  cmd.to       = fs.String("to", "", "To command to pass lines")
  cmd.verbose  = fs.Bool("v", false, "Debug info")
}

func (cmd *Run) Name() string{
  return "run"
}

func (fs *Run) Run(args []string) {
  store   := NewStore(*baseDir)
  journal := NewJournal(*fs.verbose, *baseDir + "/journal.log")
  defer store.Close()
  defer journal.Close()

  jobList := NewJobList(store, journal)

  fromListing := getFrom(*fs.from)        // result of the listing
  for _, arg := range fromListing {
    cmd := *fs.to + " " + arg
    job := NewJob(cmd, *store, *journal)
    if !jobList.Include(*job) {
      job.SetState(NEW)
      jobList.Add(*job)
    }
  }

  for _, job := range jobList {
    if job.GetState() == NEW {
      job.SetState(RUNNING)
      err := job.Run()
      if err == nil {
        job.SetState(SUCCEEDED)
      } else {
        job.SetState(FAILED)
      }
    }
  }
}

// Show
type Show struct {
  baseDir *string
  verbose *bool
  jobIds string
}

func (cmd *Show) Name() string { return "show" }

func (cmd *Show) DefineFlags(fs *flag.FlagSet) {
  cmd.baseDir  = fs.String("dir", ".tsinkf", "directory where state files are created")
  cmd.verbose = fs.Bool("v", false, "Debug info")
}

func (fs *Show) Run(jobIDs []string) {
  store   := NewStore(*baseDir)
  journal := NewJournal(*fs.verbose, *baseDir + "/journal.log")

  defer store.Close()
  defer journal.Close()

  jobList := NewJobList(store, journal)

  printable := JobSpecific(jobIDs)

  for _, job := range jobList {
    if printable(job) {
      if !*fs.verbose {
        fmt.Println(job.ToString())
      } else {
        fmt.Println(job.ToString() + "\n" + job.GetOutput())
      }
    }
  }
}

// Reset

type Reset struct {
  baseDir *string
  verbose *bool
}

func (cmd *Reset) Name() string { return "reset" }

func (cmd *Reset) DefineFlags(fs *flag.FlagSet) {
  cmd.baseDir  = fs.String("dir", ".tsinkf", "directory where state files are created")
  cmd.verbose = fs.Bool("v", false, "Debug info")
}

func (fs *Reset) Run(jobIDs []string) {
  store   := NewStore(*baseDir)
  journal := NewJournal(*fs.verbose, *baseDir + "/journal.log")

  defer journal.Close()
  defer store.Close()

  jobList   := NewJobList(store, journal)
  resetable := JobSpecific(jobIDs)

  for _, job := range jobList {
    if resetable(job) {
      job.SetState(NEW)
    }
  }
}

func main() {
  // do the no args version...
  Parse(new(Run), new(Show), new(Reset))
}
