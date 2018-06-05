package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"orfiles/subfiles"
)

type page struct {
	Msg  string
	Text string
}
type conf struct {
	Port     string `json:"port"`
	Path     string `json:"path"`
	FileType string `json:"filetype"`
}

var filesConf *subfiles.StFilesCnfg
var fileExec *subfiles.StFile
var config conf

func _checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	filesConf.CheckFiles()
	t, _ := template.ParseFiles("index.html")
	inNames := filesConf.Names[0]
	for i := 1; i <= len(filesConf.Names)-1; i++ {
		inNames = inNames + "," + filesConf.Names[i]
	}
	t.Execute(w, &page{Text: inNames})
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	delfname := r.FormValue("delf")
	filesConf.DelFile(delfname)
	http.Redirect(w, r, "/", 302)
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("New file")
	_checkErr(err)
	defer file.Close()
	filesConf.AddFile(file, handler.Filename)
	http.Redirect(w, r, "/", 302)
}
func execHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	execfname := r.FormValue("execf")
	fileExec.FileInit(execfname)
	values := fileExec.GetFileValue(filesConf.Path)
	filesConf.CheckFiles()
	inNames := filesConf.Names[0]
	for i := 1; i <= len(filesConf.Names)-1; i++ {
		inNames = inNames + "," + filesConf.Names[i]
	}
	t.Execute(w, &page{Msg: values, Text: inNames})
}

func main() {
	filesConf = new(subfiles.StFilesCnfg)
	fileExec = new(subfiles.StFile)
	byt, err := ioutil.ReadFile("./config.json")
	_checkErr(err)
	err = json.Unmarshal(byt, &config)
	_checkErr(err)
	filesConf.Path = config.Path
	filesConf.Type = config.FileType

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/exec", execHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(config.Port, nil)
}
