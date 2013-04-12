package main

import (
  "flag"
  "os"
  "os/exec"
  "log"
  "bytes"
  "strings"
  "encoding/base64"
  "fmt"
)

// Usage
// sink --from="ls -al" --to="wc -l"

//TODO: Figure out default parameters

var (
  from     = flag.String("from", "", "From command line tool")
  to       = flag.String("to", "", "To command to pass lines")
  baseDir  = flag.String("dir", "state", "directory where state files are created")
  showHelp = flag.Bool("h", false, "Show help")
)

func init() {
  flag.Parse()
}

func initState () {
  // create directory structure that will be required for writting

  for _, directory := range []string{"new", "running", "failed", "succeeded"} {
    todo := *baseDir + "/" + directory
    err := os.MkdirAll(todo, 0755)
    if err != nil {
      panic(err)
    }
  }
}

func executeCmd(targetCmd string, targetArg string) bool {
  // get
  output := targetArg + "-out.log"
  cmd := exec.Command("bash","-c",targetCmd +" &>" + output)
  err := cmd.Run()

  return err == nil
}

func getFrom(fromCmd string) []string {
  cmd := exec.Command("bash", "-c", fromCmd)
  cmd.Stdin = strings.NewReader("some input")
  var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()
  if err != nil {
    log.Fatal(err)
  }

  return strings.Split(out.String(), "\n")
}

type JobState int

const (
  UNKNOWN   JobState = -1
  NEW       JobState = 0
  RUNNING   JobState = 1
  FAILED    JobState = 2
  SUCCEEDED JobState = 3
)

type Job struct {
  hash string
  body string
  state JobState
}

func (j Job) getPath(jobState JobState) string {
  if jobState == NEW {
    return *baseDir + "/new/" + j.hash
  } else if jobState == RUNNING {
    return *baseDir + "/running/" + j.hash
  } else if jobState == FAILED {
    return *baseDir + "/failed/" + j.hash
  }

  // XXX: this may prove to be wrong
  return *baseDir + "/succeeded/" + j.hash
}

func (j Job) getFileState() (JobState, error) {
  // TODO: ensure each file occupies exactly 1 state

  if !fileExists(j.getPath(NEW)) {
    return UNKNOWN, fmt.Errorf("Job file does not exist")
  }

  if fileExists(j.getPath(RUNNING)) {
    return RUNNING, nil
  } else if fileExists(j.getPath(FAILED)) {
    return FAILED, nil
  } else if fileExists(j.getPath(SUCCEEDED)) {
    return SUCCEEDED, nil
  }

  return NEW, nil
}

type JobList map[string]*Job

func createHash(body string) string {
  msg := []byte(body)
  encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
  base64.StdEncoding.Encode(encoded, msg)
  return string(encoded)
}

func decodeHash(encoded string) string {
  decLen := base64.StdEncoding.DecodedLen(len(encoded))
  decoded := make([]byte, decLen)
  n, err := base64.StdEncoding.Decode(decoded, []byte(encoded))
  if err != nil {
    panic(err)
  }
  return string(n)
}

func fileExists(filepath string) bool {
  _, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
  if err != nil {
    if os.IsNotExist(err) {
      return false
    } else {
      panic(err)
    }
  }

  return true
}


func (jobList JobList) sync() {
  // update the state of each job based on what's in the filesystem
  // XXX: Think about inconsistencies that could occur here
  for _,job := range jobList {
    state, err := job.getFileState()
    //fmt.Printf("job state: %s\n", state)
    if err != nil {
      panic(err)
    }
    job.state = state
  }
}

func newJobList(cmdName string, listing []string) JobList {
  jobList := JobList{}
  for _, args := range listing {
    cmd := cmdName + " " + args
    hash := createHash(cmd)

    job := &Job{hash, cmd, NEW}
    jobList[hash] = job
    touchFile(job.getPath(NEW))
  }

  return jobList
}

func (t JobList) done() (doneJobs []*Job) {
  for _,job := range t {
    if job.state > RUNNING {
      doneJobs = append(doneJobs, job)
    }
  }
  return
}

func (jobList JobList) available() (availableJobs []*Job) {
  for _,job := range jobList {
    if job.state == NEW {
      availableJobs = append(availableJobs, job)
    }
  }

  return
}


func (j *Job) update(to JobState) {
  cwd, err := os.Getwd()

  if err != nil {
    panic(err)
  }

  fullPath := cwd + "/" + *baseDir
  if j.state == NEW && to == RUNNING {
    err := os.Symlink(fullPath + "/new/" + j.hash, fullPath + "/running/" + j.hash)
    if err != nil {
      panic(err)
    }
  } else if j.state == RUNNING && to == SUCCEEDED {
    err := os.Remove(fullPath + "/running/" + j.hash)
    if err != nil {
      panic(err)
    }

    err = os.Symlink(fullPath + "/new/" + j.hash, fullPath + "/succeeded/" + j.hash)
    if err != nil {
      panic(err)
    }
  } else if j.state == RUNNING && to == FAILED {
    err := os.Remove(fullPath + "/running/" + j.hash)
    if err != nil {
      panic(err)
    }

    err = os.Symlink(fullPath + "/new/" + j.hash, fullPath + "/failed/" + j.hash)
    if err != nil {
      panic(err)
    }
  }

  j.state = to
}

func touchFile (filename string) {
  _, err := os.OpenFile(filename, os.O_CREATE, 0666)
  if err != nil {
    panic(err)
  }
}

func (job *Job) run() error {
  outputFile := *baseDir + "/new/" + job.hash
  execution  := job.body + " &>" + outputFile

  //TODO setup some proper piping to clean up the process tree 
  cmd := exec.Command("bash","-c", execution)
  err := cmd.Run()

  return err
}

func main() {
  if *showHelp {
    flag.PrintDefaults()
    os.Exit(0)
  }

  initState()

  fromListing := getFrom(*from)        // result of the listing

  todo := newJobList(*to, fromListing)
  todo.sync()

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
