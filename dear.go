package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
  "log"
  "os"
  "encoding/gob"
)

type Request struct {
  Command string
  Args []string
  Origin string
  LastSender string
  Count int
}

var (
  laddr string
  peers = make(map[string]bool)
  requests = make(chan Request)
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
  enc := gob.NewEncoder(conn)
  req := Request{
    "hello",
    []string{},
    laddr,
    laddr,
    0,
  }
  enc.Encode(req)
  conn.Close()
}

func ack(addr string) {
  /* TODO: if err, don't try and check anyway */
  conn, err := net.Dial("tcp", addr)
  check("could not connect", err)

  addPeer(addr)

  enc := gob.NewEncoder(conn)
  req := Request{
    "ack",
    []string{},
    laddr,
    laddr,
    0,
  }
  enc.Encode(req)
  conn.Close()
}

func addPeer(addr string) {
  peers[addr] = true
}

func handleConnection(conn net.Conn) {
  var req Request
  dec := gob.NewDecoder(conn)
  err := dec.Decode(&req)
  if err != nil {
    log.Println(err)
    return
  }

  fmt.Printf("received %s %s\n", req.Command, req.Origin)
  switch req.Command {
    case "ack":
      addPeer(req.Origin)
    case "hello":
      ack(req.Origin)
    case "add":
      fmt.Println("I received an addition", req.Args)
      requests <- req
    case "query":
      /* TODO: also satisfy */
      fmt.Println("I received an query", req.Args)
      requests <- req
    case "get":
      /* TODO: only if not satisfied */
      fmt.Println("I received a get", req.Args)
      requests <- req
    default:
      fmt.Fprintf(conn, "wagwan")
  }

  conn.Close()
}

func sendRequests(req Request) {
  lastSender, count := req.LastSender, req.Count
  req.LastSender, req.Count = laddr, count + 1
  for peer := range(peers) {
    if (peer != req.Origin && peer != lastSender && count < 10 /*random num*/) {
      go sendRequest(req, peer)
    }
  }
}

func sendRequest(req Request, peer string) {
  conn, err := net.Dial("tcp", peer)
  check("Couldn't connect to peer", err)
  enc := gob.NewEncoder(conn)
  enc.Encode(req)
  conn.Close()
}

func makeRequest(cmd string, args []string) Request {
  return Request{
    cmd,
    args,
    laddr,
    laddr,
    0,
  }
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
        requests <- makeRequest(cmd, args)
      case "query":
        requests <- makeRequest(cmd, args)
      case "get":
        requests <- makeRequest(cmd, args)
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

func broadcastRequests() {
  for {
    req := <-requests
    sendRequests(req)
  }
}

func getIP() string{
   addrs, err := net.InterfaceAddrs()
   if err != nil {
     panic(err)
   }
    for _, addr := range addrs {
      ip := strings.Split(addr.String(), "/")[0]
      nums := strings.Split(ip, ".")
      if (nums[0] == "192" && nums[1] == "168") {
        return ip
     }
   }
   return "localhost"
 }

func main() {

  listener, err := net.Listen("tcp", ":0")
  if err != nil {
    panic(err)
  }
  laddr = getIP()
  fmt.Printf("listening on %s\n", laddr)
  stdin := bufio.NewScanner(os.Stdin)

  go readCommands(stdin)
  go broadcastRequests()

  for {
    conn, err := listener.Accept()
    if err != nil {
      log.Println("error accepting connection")
    }
    go handleConnection(conn)
  }
}

