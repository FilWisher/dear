package main

import (
  "net/http"
  "fmt"
  "io/ioutil"
  "github.com/julienschmidt/httprouter"
  "github.com/filwisher/digestif"
  "text/template"
)

type Item struct {
  Name string
  Index int
}

type Content struct {
  Name string
}

func Index(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  file, _ := ioutil.ReadFile("index.html")
  res.Header().Add("content-type", "text/html")
  fmt.Fprintf(res, string(file))
}

func List(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  item := template.Must(template.ParseFiles("item.html"))
  files, _ := ioutil.ReadDir("data/")
  fmt.Fprintf(res, "<!DOCTYPE html><html><head></head><body>")
  for i, file := range files {
    it := Item{
      file.Name(),
      i+1,
    }
    item.Execute(res, it)
  }
  fmt.Fprintf(res, "</body></html>")
}

func Save(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  req.ParseForm()
  name := req.Form["name"][0]
  res.Header().Add("content-type", "text/plain")
  data := []byte("Hello there " + name)
  filename := digest.ToHexString(digest.Hash(data))
  digest.Save("data/" + filename, data)
  http.Redirect(res, req, "/hello/" + filename, 302)
}

func Show(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
  item := template.Must(template.ParseFiles("content.html"))
  file := params.ByName("name")
  data, _ := ioutil.ReadFile("data/" + file)
  cont := Content{
    string(data),
  }
  item.Execute(res, cont)
}

func main() {
  router := httprouter.New()
  router.GET("/", Index)
  router.GET("/ls", List)
  router.POST("/hello", Save)
  router.GET("/hello/:name", Show)
  http.ListenAndServe(":8080", router)
  fmt.Println("Hello")
}
