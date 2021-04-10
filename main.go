package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

func main() {
	tree := GeneratePagesTree("pages")

	HandleFunc("/about", aboutEndpointHandler)

	HandleTree(&tree)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	v := reflect.ValueOf(http.DefaultServeMux).Elem()
	fmt.Printf("routes: %v\n", v.FieldByName("m"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
func aboutEndpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the about page, ")
	fmt.Fprintf(w, "where you can find information about us.")
}
