package main

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type post struct {
	file  string
	route string

}
func generate(fileName string) post {
	return post{
		file:  "./posts/" + fileName,
		route: "/posts/" + strings.TrimSuffix(fileName, filepath.Ext(fileName)),
	}
}

func main() {
	registerPosts()

	http.HandleFunc("/ping", ping)
	http.HandleFunc(
		"/main", func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadFile("index.html")
			if err != nil {
				fmt.Fprintln(w, "File reading error", err)
				return
			}
			w.Write(data)
		},
	)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerPosts() {
	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		p := generate(file.Name())
		fmt.Println(p)
		http.HandleFunc(
			p.route, func(w http.ResponseWriter, r *http.Request) {
				data, err := ioutil.ReadFile(p.file)
				if err != nil {
					fmt.Fprintln(w, "File reading error", err)
					return
				}
				w.Write(blackfriday.Run(data))
			},
		)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
