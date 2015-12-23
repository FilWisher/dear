package main

import (
  "net"
  "fmt"
  "os"
  "github.com/filwisher/digestif"
  "io/ioutil"
  "bufio"
)

const PREFIX string = "data"

func handleConnection(conn *net.Conn) {
  fmt.Println("Got a connection")
  fmt.Fprintf(*conn, "Just got this")
}

func getDescription() string {
  reader := bufio.NewReader(os.Stdin)
  fmt.Println("Please enter a description:")
  text, _ := reader.ReadString('\n')
  return text
}

func addToIndex(hash, filename string) {

  description := getDescription()

  if _, err := os.Stat(PREFIX + "/index"); os.IsNotExist(err) {
    os.Create(PREFIX + "/index")
  }
  f, err := os.OpenFile(PREFIX + "/index", os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  _, err = f.WriteString(hash + "\t" + filename + "\t" + description)

  if err != nil {
    panic(err)
  }
}

func save() {
  if len(os.Args) < 3 {
    panic("not enough arguments")
  }
  if _, err := os.Stat(PREFIX); os.IsNotExist(err) {
    os.Mkdir(PREFIX, 0700)
  }
  data, err := ioutil.ReadFile(os.Args[2])
  if err != nil {
    panic("could not open file")
  }
  filename := digest.ToHexString(digest.Hash(data))
  addToIndex(filename, os.Args[2])
  digest.Save(PREFIX + "/" + filename, data)
}

func list() {
  data, err := ioutil.ReadFile(PREFIX + "/index")
  if err != nil {
    panic(err)
  }
  fmt.Println(string(data))
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
    case "save":
      save()
    case "list":
      list()
    default:
      fmt.Println("Dont recognize that comand")
  }
}
