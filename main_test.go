package main

import (
  "testing"
  "bytes"
  //"strings"
)


func TestRawDetails(t *testing.T) {
  details := getRowDetails(0, true)
  if details != "\033[1m\033[4m         1\033[0m\033[0m|" {
    t.Error("Detail is wrong [" + details + "] should be [\033[1m\033[4m         1|\033[0m\033[0m]")
  }

  details = getRowDetails(1, false)
  if details != "         2|" {
    t.Error("Detail is wrong [" + details + "]")
  }

}

func BenchmarkDrawRows(b *testing.B) {
  str := randBook(7000)
  ce := NewCatpadEditor(10, 90, 1, str)
  buff := bytes.NewBufferString("")

  for n:=0; n<b.N; n++ {
    editorDrawRows(&ce, buff)
    buff.Reset()
  }
}
