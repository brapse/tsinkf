package main

import (
  "crypto/md5"
  "io"
  "encoding/hex"
)

func CreateHash(cmd string) string {
  // cmd5 the command
  h := md5.New()
  io.WriteString(h, cmd)

  return hex.EncodeToString(h.Sum(nil))
}
