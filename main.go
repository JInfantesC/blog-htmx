package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	pagesDir, err := ReadDirectory("pages")
	if err != nil {
		log.Panicln(err)
	}
	verbosePrintPage(os.Stdout, &pagesDir)

	HandlePage(
		&Page{
			webPath:      "/",
			isType:       PAGE,
			bufferedData: loadTemplate("templates/index.gohtml", pagesDir),
		},
	)
	HandleDirectory(&pagesDir)

	fs := http.FileServer(http.Dir("static/"))
	Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loadTemplate(route string, data interface{}) []byte {
	t, err := template.ParseFiles(route)
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{} // Implementa io.Write
	err = t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// verbosePrintPage recorre Pages invocando String() de cada p y sus p.SubPages
func verbosePrintPage(wr io.Writer, p *Page) {
	fmt.Fprintln(wr, p.String())
	for _, subPage := range p.SubPages {
		verbosePrintPage(wr, &subPage)
	}
}
