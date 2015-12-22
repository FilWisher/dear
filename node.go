package main

/*
    TODO:
    Define procedures that are callable on peers
*/

import "fmt"

type Command struct {
  Name string
}

type Response struct {
  Body string
}

type Node struct {}

func (n *Node) Run (com Command, res *Response) error {
  fmt.Println("Received call to run")
  res.Body = "My name is Wil\n"
  return nil
}
