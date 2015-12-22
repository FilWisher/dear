package main

import (
  "net"
  "net/rpc"
  "fmt"
  "os"
  "syscall"
  "os/signal"
)

func wait() {
  signals := make(chan os.Signal)
  signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
  <-signals
}

func server() {

  node := new(Node)
  rpc.Register(node)

  s, err := net.Listen("tcp", ":8080")
  if err != nil {
    panic("net listen")
  }

  fmt.Println("server")

  go rpc.Accept(s)

  wait()
}
