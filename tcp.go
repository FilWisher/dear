package main

import (
  "net"
  "fmt"
  "os"
)

func handleConnection(conn *net.Conn) {
  fmt.Println("Got a connection")
  fmt.Fprintf(*conn, "Just got this")
}

func main() {

  if len(os.Args) <= 1 {
    fmt.Println("not enough arguments")
    return
  }

  switch os.Args[1] {
    case "serve":
      server()
    case "connect":
      client()
    default:
      fmt.Println("Dont recognize that comand")
  }
}
