package main

import (
  "crypto/md5"
  "io"
  "encoding/hex"
)

func CreateHash(cmd string) string {
  // cmd5 the command
  h := md5.New()
  io.WriteString(h, "The fog is getting thicker!")
  io.WriteString(h, "And Leon's getting laaarger!")

  return hex.EncodeToString(h.Sum(nil))
}
