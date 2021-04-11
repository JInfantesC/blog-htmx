package main

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type PageType uint8

const (
	PAGE PageType = iota
	DIRECTORY
)

/*Page baseType*/
type page struct {
	route   string
	isType  PageType
	base    string
	webName string
	data    []byte
	id      string
}

func (p page) WebName() string {
	return p.webName
}
func (p page) Base() string {
	return p.base
}
func (p page) Route() string {
	return p.route
}
func (p page) Type() PageType {
	return p.isType
}
func (p *page) HtmlId() string {
	if p.id == "" {
		p.id = strconv.Itoa(int(rand.Uint32()))
	}
	return p.id
}

type pageFile struct {
	page
	extension string
}

type pageDir struct {
	page
	dirList   []pageDir
	pagesList []pageFile
}

func (p pageDir) PagesList() []pageFile {
	return p.pagesList
}
func (p pageDir) DirList() []pageDir {
	return p.dirList
}

func GeneratePagesTree(dirname string) pageDir {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	dir := pageDir{
		page: page{
			route:  dirname,
			isType: DIRECTORY,
			base:   filepath.Base(dirname),
		},
		dirList:   make([]pageDir, 0),
		pagesList: make([]pageFile, 0),
	}
	for _, file := range files {

		if file.IsDir() {
			dir.dirList = append(
				dir.dirList, GeneratePagesTree(filepath.Join(dirname, file.Name())),
			)

		} else {
			dir.pagesList = append(
				dir.pagesList, GeneratePage(file, filepath.Join(dirname, file.Name())),
			)
		}
	}

	return dir
}

func GeneratePage(f fs.FileInfo, dirname string) pageFile {
	return pageFile{
		page: page{
			route:   dirname,
			isType:  PAGE,
			base:    filepath.Base(dirname),
			webName: strings.TrimSuffix("/"+dirname, filepath.Ext(f.Name())),
		},
		extension: filepath.Ext(f.Name()),
	}
}

func HandleTree(dir *pageDir) {
	HandleDir(dir)
	for i := 0; i < len(dir.dirList); i++ {
		HandleTree(&dir.dirList[i])
	}

	for i := 0; i < len(dir.pagesList); i++ {
		// Con este bucle for, para recuperar la dirección de memoria correcta de pagesList[i],
		//con el anterior reúsa &p y en pagesList con más de una pageFile siempre devuelve el último
		HandlePage(&dir.pagesList[i])
	}
}

func HandleDir(dir *pageDir) {
	HandleFunc(
		"/"+dir.route+"/", func(w http.ResponseWriter, r *http.Request) {
			if dir.data == nil {
				t, err := template.New("list").Parse(`<section class="content"><h3>{{ .Base }}</h3>
					{{range .DirList}}
						<button hx-post="{{ .Route }}"
								hx-trigger="click"
								hx-target="#content_{{ .HtmlId }}"
								hx-swap="innerHTML">
							{{.Base}}
						</button>
						<article id="content_{{ .HtmlId }}" class="card"></article>
					{{end}}
					{{range .PagesList}}
						<button hx-post="{{ .WebName }}"
								hx-trigger="click"
								hx-target="#content_{{ .HtmlId }}"
								hx-swap="innerHTML">
							{{.Base}}
						</button>
						<article id="content_{{ .HtmlId }}" class="card"></article>
					{{end}}
					</section>`)
				if err != nil {
					panic(err)
				}

				err = t.Execute(w, dir)
				if err != nil {
					panic(err)
				}
			} else {
				w.Write(dir.data)
			}

		},
	)
}

func HandlePage(p *pageFile) {
	HandleFunc(
		p.webName, func(w http.ResponseWriter, r *http.Request) {
			if p.data == nil {
				data, err := ioutil.ReadFile(p.Route())
				if err != nil {
					fmt.Fprintln(w, "File reading error", err)
					return
				}
				if p.extension == ".md" {
					p.data = blackfriday.Run(data)
				} else {
					p.data = data
				}
			}
			w.Write(p.data)
		},
	)
}
