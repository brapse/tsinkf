package main

import (
  "strings"
  "fmt"
  "os"
  "time"
)

type Journal struct {
  *os.File
  writeStdout bool
}

func NewJournal(stdout bool, fileLoc string) *Journal {
  fp, err := os.OpenFile(fileLoc, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
  if err != nil {
    panic(err)
  }

  return  &Journal{fp, stdout}
}

func (j Journal) Log(job Job, toState JobState) {
    fromLabel := STATELABELS[job.GetState()]
    toLabel   := STATELABELS[toState]
    msg := strings.Join([]string{
      FormatTime(time.Now()),
      fromLabel + "->" + toLabel,
      job.cmd,
      job.id}, "\t") + "\n"

    if j.writeStdout {
        fmt.Printf(msg)
    }
    if _, err := j.WriteString(msg); err != nil {
        panic(err)
    }
}
