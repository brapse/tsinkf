package main

import (
  "strings"
  "fmt"
  "os"
  "time"
)

type Journal struct {
  writeStdout bool
  fp *os.File
}

func NewJournal(stdout bool, fileLoc string) *Journal {
  fp, err := os.OpenFile(fileLoc, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
  if err != nil {
    panic(err)
  }

  return  &Journal{stdout, fp}
}

func (j Journal) log(job Job, toState JobState) {
    fromLabel := STATELABELS[job.state]
    toLabel   := STATELABELS[toState]
    msg := strings.Join([]string{
      time.Now().Format("2006-01-02 15:04:05"),
      fromLabel + "->" + toLabel,
      job.cmd,
      job.hash }, "\t") + "\n"

    if j.writeStdout {
        fmt.Printf(msg)
    }
    if _, err := j.fp.WriteString(msg); err != nil {
        panic(err)
    }
}

func (j Journal) close() {
    err := j.fp.Close()
    if err != nil {
        panic(err)
    }
}
