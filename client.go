package main

import (
  "net/rpc"
  "fmt"
  "io"
)

/* 
    TODO: 
    Allow different remote procedures to be callable
    Introduce peer-to-peer logic 
*/

/*
func readResponse(res io.Reader) string {
  data := make([]byte, 512)
  length, err := res.Read(data)
  if err != nil {
    panic("response not right")
  }
  return string(data[:length])
}
*/

func client() {

  cli, err := rpc.Dial("tcp", "localhost:8080")
  if err != nil {
    panic("cannot dial")
  }

  command := &Command{"My name is Wil"}
  var res Response
  err = cli.Call("Node.Run", command, &res)
  if err != nil {
    panic("could'n call command")
  }
  fmt.Printf("Got: %s", res.Body)
  cli.Close()
}
