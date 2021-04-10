package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

func main() {
	tree := GeneratePagesTree("pages")

	HandlePage(&pageFile{
		page: page{
			route:   "index.html",
			isType:  PAGE,
			webName: "/",
			data:    nil,
		},
		extension: ".html",
	})
	HandleTree(&tree)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	v := reflect.ValueOf(http.DefaultServeMux).Elem()
	fmt.Printf("routes: %v\n", v.FieldByName("m"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
