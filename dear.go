package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
  "log"
  "os"
)

var (
  laddr string
  peers = make(map[string]bool)
)

func check(msg string, err error) {
  if err != nil {
    log.Println(msg, err)
  }
}

/* ask to join with a node */
func greet(addr string) {
  /* TODO: if err, don't try and check anyway */
  conn, err := net.Dial("tcp", addr)
  check("could not connect", err)
  fmt.Fprintf(conn, "hello %s\n", laddr)
  conn.Close()
}

func ack(addr string) {
  /* TODO: if err, don't try and check anyway */
  conn, err := net.Dial("tcp", addr)
  check("could not connect", err)
  fmt.Fprintf(conn, "ack %s\n", laddr)
  addPeer(addr)
  conn.Close()
}

func addPeer(addr string) {
  peers[addr] = true
}

func handleConnection(conn net.Conn) {
  /* TODO: don't presume everything received is ack */
  scanner := bufio.NewScanner(conn)
  scanner.Scan()
  /* TODO: check error */
  input := strings.Split(scanner.Text(), " ")
  cmd, args := input[0], input[1:]
  fmt.Printf("received %s %s\n", cmd, args[0])
  addPeer(args[0])

  conn.Close()
}

func readCommands(scanner *bufio.Scanner) {
  for scanner.Scan() {
    if err := scanner.Err(); err != nil {
      log.Println("error reading commmands", err)
    }
    input := strings.Split(scanner.Text(), " ")
    cmd, args := input[0], input[1:]
    switch cmd {
      case "add":
        fmt.Println("add", args)
      case "query":
        fmt.Println("query", args)
      case "get":
        fmt.Println("get", args)
      case "greet":
        greet(args[0])
      case "peers":
        for peer := range(peers) {
          fmt.Println(peer)
        }
      case "help":
        fmt.Println("help: get, query, or add")
      default:
        fmt.Println("command not known: try help")
    }
  }
}

func main() {

  listener, err := net.Listen("tcp", ":0")
  if err != nil {
    panic(err)
  }
  laddr = listener.Addr().String()
  fmt.Printf("listening on %s\n", laddr)
  stdin := bufio.NewScanner(os.Stdin)

  go readCommands(stdin)

  for {
    conn, err := listener.Accept()
    if err != nil {
      log.Println("error accepting connection")
    }
    go handleConnection(conn)
  }
}

