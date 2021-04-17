package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type PageType uint8

func (pt PageType) String() string {
	switch pt {
	case PAGE:
		return "Page"
	case DIRECTORY:
		return "Dir"
	}
	return "Error"
}

const (
	PAGE PageType = iota
	DIRECTORY
)

/*Page baseType*/
type Page struct {
	webPath      string   // Ruta o pattern donde se ejecutará en el servidor
	filePath     string   // Ruta física en el sistema de archivos
	isType       PageType // Tipo de Page
	bufferedData []byte   // Tras primera lectura. []byte a enviar por el servidor
	subPages     []Page   // Listado de sub páginas de este Page
}

func (p Page) GetExtension() string {
	if p.GetType() == PAGE {
		return filepath.Ext(p.GetFileName())
	}
	return ""
}
func (p Page) GetFileName() string {
	return filepath.Base(p.filePath)
}
func (p Page) GetType() PageType {
	return p.isType
}
func (p *Page) String() string {
	return fmt.Sprintf(
		"%s:\twebPath (filePath)-> %s (%s)", p.GetType(), p.webPath, p.filePath,
	)
}

func (p *Page) GetSubDirectories() ([]Page, error) {
	if p.isType == DIRECTORY {
		var returnPages []Page
		for _, subPage := range p.subPages {
			if subPage.isType == DIRECTORY {
				returnPages = append(returnPages, subPage)
			}
		}
		return returnPages, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%s no tiene subDirectorios", p.filePath))
	}
}
func (p *Page) GetSubPages() ([]Page, error) {
	if p.isType == DIRECTORY {
		var returnPages []Page
		for _, subPage := range p.subPages {
			if subPage.isType == PAGE {
				returnPages = append(returnPages, subPage)
			}
		}
		return returnPages, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%s no tiene subPáginas", p.filePath))
	}
}

// ReadDirectory recibe la ruta al directorio relativa a working directory y devuelve un Page
//completo listo para servir en un HandleFunc
func ReadDirectory(dirName string) (Page, error) {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return Page{}, err
	}

	dir := Page{
		webPath:  dirName,
		filePath: dirName,
		isType:   DIRECTORY,
	}

	for _, file := range files {
		if file.IsDir() {
			subDir, errDir := ReadDirectory(filepath.Join(dir.filePath, file.Name()))
			if errDir != nil {
				return Page{}, errDir
			}
			dir.subPages = append(dir.subPages, subDir)
		} else {
			subPage, errPage := ReadPageFile(file, filepath.Join(dir.filePath, file.Name()))
			if errPage != nil {
				return Page{}, errPage
			}
			dir.subPages = append(dir.subPages, subPage)
		}
	}
	return dir, nil
}

func ReadPageFile(f fs.FileInfo, dirName string) (Page, error) {
	if f.IsDir() {
		return Page{}, errors.New(
			fmt.Sprintf(
				"%s es un directorio, ReadPageFile requiere un archivo",
				f.Name(),
			),
		)
	}

	return Page{
		webPath:  strings.TrimSuffix(dirName, filepath.Ext(dirName)),
		filePath: dirName,
		isType:   PAGE,
	}, nil
}

/*
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
				t, err := template.New("list").Parse(
					`<section class="content"><h3>{{ .Base }}</h3>
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
					</section>`,
				)
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
*/
