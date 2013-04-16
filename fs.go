package main

import "os"

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

