package main

import (
  "encoding/base64"
  "strings"
)

func CreateHash(cmd string) string {
  msg := []byte(cmd)
  encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
  base64.StdEncoding.Encode(encoded, msg)
  return string(encoded)
}

var null = string([]byte{0})
func DecodeHash(encoded string) string {
  decLen := base64.StdEncoding.DecodedLen(len(encoded))
  decoded := make([]byte, decLen)
  _, err := base64.StdEncoding.Decode(decoded, []byte(encoded))
  if err != nil {
    panic(err)
  }

  // XXX: Remove null byte(s)
  return strings.Replace(string(decoded), null, "", -1)
}
