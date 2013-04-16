package main

import (
  "flag"
  "os"
  "os/exec"
  "strings"
  "bytes"
  "log"
)


// TODO
// + subcommands: run, show, reset
// + Refactor
// + Refactor Job stuff
// + Refactor execution stuff
// + Add journal style logging

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

func getFrom(fromCmd string) []string {
  cmd := exec.Command("bash", "-c", fromCmd)
  var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()
  if err != nil {
    log.Fatal(err)
  }
  return strings.Split(out.String(), "\n")
}

func main() {
  if *showHelp {
    flag.PrintDefaults()
    os.Exit(0)
  }

  // Create a Filestore
  store := NewStore(*baseDir)
  journal := NewJournal(*debug, *baseDir + "/journal.log")
  // Create a Journal

  // RUNNING
  fromListing := getFrom(*from)        // result of the listing
  jobList     := NewJobList(*to, fromListing, *store, *journal)

  for _, job := range jobList.available() {
    jobList.update(job, RUNNING)
    err := job.run()
    if err == nil {
      jobList.update(job, SUCCEEDED)
    } else {
      jobList.update(job, FAILED)
    }
  }
}
