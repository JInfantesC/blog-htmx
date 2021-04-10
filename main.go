package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func main() {
	tree := GeneratePagesTree("pages")

	HandleTree(&tree)

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
	v := reflect.ValueOf(http.DefaultServeMux).Elem()
	fmt.Printf("routes: %v\n", v.FieldByName("m"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
