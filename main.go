package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

func main() {
	pageDirectory := getPageDirectoryFlag()

	pagesDir, err := ReadDirectory(pageDirectory)
	if err != nil {
		log.Panicln(err)
	}
	verbosePrintPage(os.Stdout, &pagesDir)

	HandlePage(
		&Page{
			WebPath:      "/",
			isType:       PAGE,
			bufferedData: loadTemplate("templates/index-htmx.gohtml", pagesDir),
		},
	)
	HandleDirectory(&pagesDir)

	Handle("/static/", http.FileServer(http.FS(static)))

	log.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loadTemplate(route string, data interface{}) []byte {
	//t, err := template.ParseFiles(route)
	t, err := template.ParseFS(templates, route)
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

func getPageDirectoryFlag() string {
	var dir string
	flag.StringVar(&dir, "dir", "pages", "Dirección del directorio a servir")

	flag.Parse()
	return dir
}
