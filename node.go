/*
    Defines the API for remote calls to Nodes
*/
package main

import (
  "io/ioutil"
)

type Command struct {
  Name string
}

type Response struct {
  Body string
}

type Node struct {}

/*
    List files available at target node. Available files
    returned in Response.Body
*/
func (n *Node) List (com Command, res *Response) error {
  data, err := ioutil.ReadFile(PREFIX + "/index")
  if err != nil {
    return err
  }
  res.Body = string(data)
  return nil
}

/*
    Get a file from target node. Filename (hash) passed in
    Command.Name. File returned in Response.Body
*/
func (n *Node) Get (com Command, res *Response) error {
  data, err := ioutil.ReadFile(PREFIX + "/" + com.Name)
  if err != nil {
    return err
  }
  res.Body = string(data)
  return nil
}
