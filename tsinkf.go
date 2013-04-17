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


// TODO
// + subcommands: run, show, reset
// + Refactor Job stuff [DONE]
// + Refactor execution stuff
// + Add journal style logging [DONE]

var (
  from     = flag.String("from", "", "From command line tool")
  to       = flag.String("to", "", "To command to pass lines")
  baseDir  = flag.String("dir", ".tsinkf", "directory where state files are created")
  showHelp = flag.Bool("h", false, "Show help")
  debug    = flag.Bool("v", false, "Debug info")
)

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
  baseDir *string
  debug *bool
}

func (cmd *Run) DefineFlags(fs *flag.FlagSet) {
  cmd.from     = fs.String("from", "", "From command line tool")
  cmd.to       = fs.String("to", "", "To command to pass lines")
  cmd.baseDir  = fs.String("dir", ".tsinkf", "directory where state files are created")
  cmd.debug    = fs.Bool("v", false, "Debug info")
}

func (cmd *Run) Name() string{
  return "run"
}

func (fs *Run) Run() {
  store   := NewStore(*fs.baseDir)
  journal := NewJournal(*fs.debug, *fs.baseDir + "/journal.log")
  defer store.Close()
  defer journal.Close()

  jobList := NewJobList(store, journal)

  fromListing := getFrom(*fs.from)        // result of the listing
  for _, arg := range fromListing {
    cmd := *to + " " + arg
    job := NewJob(cmd)
    if !jobList.include(job) {
      job.state = NEW
      jobList.add(*job)
    }
  }

  for _, job := range jobList {
    if job.state == NEW {
      jobList.update(&job, RUNNING)
      err := job.run()
      if err == nil {
        jobList.update(&job, SUCCEEDED)
      } else {
        jobList.update(&job, FAILED)
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
  store := NewStore(*fs.baseDir)
  journal := NewJournal(*fs.verbose, *fs.baseDir + "/journal.log")

  defer store.Close()
  defer journal.Close()

  jobList := NewJobList(store, journal)

  for _, job := range jobList {
    fmt.Println(job.toString())
  }
}

func main() {
  // do the no args version...
  switch os.Args[1] {
    case "run":
      Parse(new(Run))
    case "show":
      Parse(new(Show))
    default:
      fmt.Printf("invalid command %s", os.Args[1])
  }
}
