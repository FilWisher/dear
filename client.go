package main

import (
  "net/rpc"
  "fmt"
  "os"
)

/* 
    TODO: 
    Introduce peer-to-peer logic 
*/

func listRemote(client *rpc.Client) {
  var res Response
  command := &Command{}
  err := client.Call("Node.List", command, &res)
  if err != nil {
    panic(err)
  }
  fmt.Println(res.Body)
}

func getRemote(client *rpc.Client) {
  var res Response
  if len(os.Args) < 4 {
    fmt.Println("Not enough arguments :(")
  }
  command := &Command{os.Args[3]}
  err := client.Call("Node.Get", command, &res)
  if err != nil {
    panic(err)
  }
  fmt.Println(res.Body)
}

func client() {

  cli, err := rpc.Dial("tcp", "localhost:8080")
  if err != nil {
    panic("cannot dial")
  }
  if len(os.Args) < 3 {
    panic("not enough arguments")
  }
  switch os.Args[2] {
    case "list":
      listRemote(cli)
    case "get":
      getRemote(cli)
    default:
      fmt.Println("Unrecognized command")
      fmt.Println(os.Args[2])
  }

  cli.Close()
}
