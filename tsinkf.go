package main

import (
  "flag"
  "os"
  "jobs"
  "strings"
  "exec"
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

  // RUNNING
  fromListing := getFrom(*from)        // result of the listing
  jobList     := NewJobList(*to, fromListing)

  for _, job := range todo.available() {
    job.update(RUNNING)
    err := job.run()
    if err == nil {
      job.update(SUCCEEDED)
    } else {
      job.update(FAILED)
    }
  }
}
