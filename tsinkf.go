package main

import (
  "flag"
  "os"
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
  debug *bool
}

func (cmd *Run) DefineFlags(fs *flag.FlagSet) {
  cmd.from     = fs.String("from", "", "From command line tool")
  cmd.to       = fs.String("to", "", "To command to pass lines")
  cmd.debug    = fs.Bool("v", false, "Debug info")
}

func (cmd *Run) Name() string{
  return "run"
}

func (fs *Run) Run() {
  store   := NewStore(*baseDir)
  journal := NewJournal(*fs.debug, *baseDir + "/journal.log")
  defer store.Close()
  defer journal.Close()

  jobList := NewJobList(store, journal)

  fromListing := getFrom(*fs.from)        // result of the listing
  for _, arg := range fromListing {
    cmd := *fs.to + " " + arg
    job := NewJob(cmd)
    if !jobList.Include(job) {
      job.state = NEW
      jobList.Add(*job)
    }
  }

  for _, job := range jobList {
    if job.state == NEW {
      jobList.Update(&job, RUNNING)
      err := job.Run()
      if err == nil {
        jobList.Update(&job, SUCCEEDED)
      } else {
        jobList.Update(&job, FAILED)
      }
    }
  }
}

// Show
type Show struct {
  baseDir *string
  verbose *bool
  jobs []string
}

func (cmd *Show) Name() string { return "show" }

func (cmd *Show) DefineFlags(fs *flag.FlagSet) {
  cmd.baseDir  = fs.String("dir", ".tsinkf", "directory where state files are created")
  cmd.verbose = fs.Bool("v", false, "Debug info")
  cmd.jobs    = fs.Args()
}

func (fs *Show) Run() {
  store   := NewStore(*baseDir)
  journal := NewJournal(*fs.verbose, *baseDir + "/journal.log")

  defer store.Close()
  defer journal.Close()

  jobList := NewJobList(store, journal)

  for _, job := range jobList {
    if !*fs.verbose {
      fmt.Println(job.ToString())
    } else {
      fmt.Println(job.ToString() + "\n" + job.Content())
    }
  }
}

// Reset

type Reset struct {
  baseDir *string
}

func (cmd *Reset) Name() string { return "reset" }

func (cmd *Reset) DefineFlags(fs *flag.FlagSet) {
}

func (fs *Reset) Run() {
  store   := NewStore(*baseDir)

  store.Reset()
  store.Close()
}

func main() {
  // do the no args version...
  switch os.Args[1] {
    case "run":
      Parse(new(Run))
    case "show":
      Parse(new(Show))
    case "reset":
      Parse(new(Reset))
    default:
      fmt.Printf("invalid command %s", os.Args[1])
  }
}
