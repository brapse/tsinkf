package main

import (
	"os"
	"io/ioutil"
  "time"
)

func TouchFile (filename string) {
  _, err := os.OpenFile(filename, os.O_CREATE, 0666)
  if err != nil {
    panic(err)
  }
}

func FileExists(filepath string) bool {
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

func RemoveFile(filepath string) {
  if FileExists(filepath) {
    err := os.Remove(filepath)
    if err != nil && os.IsNotExist(err) {
      panic(err)
    }
  }
}

func ListFiles(basePath string) (res []string) {
	fileInfos, err := ioutil.ReadDir(basePath)
	if err != nil {
		panic(err)
	}

	for _, fi := range fileInfos {
		res = append(res, fi.Name())
	}

	return res
}

func ReadFile(filepath string) string {
  content, err := ioutil.ReadFile(filepath)

  if err != nil {
    panic(err)
  }

  return string(content)
}

func LastTouch(filepath string) time.Time{
  info, err := os.Stat(filepath)

  if err != nil {
    panic(err)

  }

  return info.ModTime()
}