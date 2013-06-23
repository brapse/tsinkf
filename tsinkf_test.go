package main

import (
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

var CMD_SUCCESS = 0
var CMD_FAILURE = 1

func matches(needle string, heystack string) bool {
	match, err := regexp.MatchString(needle, heystack)
	if err != nil {
		panic(err)
	}

	return match
}

func resetState() {
	cmd := exec.Command("rm", "-rf", ".tsinkf")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func tsinkfExec(args string) (string, int) {
	cmdLine := []string{"go",
                      "run",
                      "utils.go",
                      "jobs.go",
                      "journal.go",
                      "store.go",
                      "hashing.go",
                      "fs.go",
                      "tsinkf.go",
                      "cmd.go",
                      "run.go",
                      "show.go",
                      "reset.go"}

	cmdLine = append(cmdLine, args)

	cmd := exec.Command("bash", "-c", strings.Join(cmdLine, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), CMD_FAILURE
	}

	return string(output), CMD_SUCCESS
}

func TestRun(t *testing.T) {
	resetState()

	cmd := "run echo OKOKOK"
	output, status := tsinkfExec(cmd)

	if status != CMD_SUCCESS {
		t.Fatal("Failed to execute: ", cmd, "\n", output)
	}

	if !matches("OKOKOK", output) {
		t.Fatal("Running did not produce the expected output \"OKOKOK\"")
	}

	output, status = tsinkfExec(cmd)

	if status != CMD_SUCCESS {
		t.Fatal("Failed to execute a second time: ", cmd)
	}

	if !matches("^$", output) {
		t.Fatal("Re-run should not re-exute but somehow produced output!\n", output)
	}

	cmd = "show -v"
	output, status = tsinkfExec(cmd)

	if status != CMD_SUCCESS {
		t.Fatal("Failed to execute a show: ", cmd)
	}

	if !matches("SUCCEEDED", output) {
		t.Fatal("Show log successful statae!\n", output)
	}

	if !matches("OKOKOK\n", output) {
		t.Fatal("Should include command output\n", output)
	}
}
