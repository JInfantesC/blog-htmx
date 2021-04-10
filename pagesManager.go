package main

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type PageType uint8

const (
	PAGE PageType = iota
	DIRECTORY
)

/*Page baseType*/
type page struct {
	route  string
	isType PageType
}

func (p page) GetRoute() string {
	return p.route
}
func (p page) GetType() PageType {
	return p.isType
}

type pageFile struct {
	page
	webName   string
	extension string
	data      []byte
}

type pageDir struct {
	page
	dirList   []pageDir
	pagesList []pageFile
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
			route:  dirname,
			isType: PAGE,
		},
		webName:   strings.TrimSuffix("/"+dirname, filepath.Ext(f.Name())),
		extension: filepath.Ext(f.Name()),
		data:      nil,
	}
}

func HandleTree(dir *pageDir) {
	for _, d := range dir.dirList {
		HandleTree(&d)
	}
	/*for _, p := range dir.pagesList {
		HandleTree(&p)
	}*/
	for i := 0; i < len(dir.pagesList); i++ {
		// Con este bucle for, para recuperar la dirección de memoria correcta de pagesList[i],
		//con el anterior reúsa &p y en pagesList con más de una pageFile siempre devuelve el último
		HandlePage(&dir.pagesList[i])
	}
}

func HandlePage(p *pageFile) {
	http.HandleFunc(
		p.webName, func(w http.ResponseWriter, r *http.Request) {
			if p.data == nil {
				data, err := ioutil.ReadFile(p.GetRoute())
				if err != nil {
					fmt.Fprintln(w, "File reading error", err)
					return
				}
				if p.extension == "md" {
					p.data = blackfriday.Run(data)
				} else {
					p.data = data
				}
			}
			w.Write(p.data)
		})
}
