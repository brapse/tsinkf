package main

import "time"

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

type action func(job Job) bool

// JobSpecific is a function that takes  a set of jobs
// and returns a function that takes a job and returns
// true if the job is member of the jobList and false otherwise
func JobSpecific(jobIDs []string) action {
	return func(target Job) bool {
		if len(jobIDs) > 0 {
			for _, jobID := range jobIDs {
				if jobID == target.id {
					return true
				}
			}
			return false
		}
		return true
	}
}
