package main

import (
	"errors"
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
	WebPath      string   // Ruta o pattern donde se ejecutará en el servidor
	filePath     string   // Ruta física en el sistema de archivos
	isType       PageType // Tipo de Page
	bufferedData []byte   // Tras primera lectura. []byte a enviar por el servidor
	SubPages     []Page   // Listado de sub páginas de este Page
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
func (p Page) String() string {
	return fmt.Sprintf(
		"%s\t%s\t%s", p.GetType(), p.WebPath, p.filePath,
	)
}

func (p *Page) GetSubDirectories() []Page {
	if p.isType == DIRECTORY {
		var returnPages []Page
		for _, subPage := range p.SubPages {
			if subPage.isType == DIRECTORY {
				returnPages = append(returnPages, subPage)
			}
		}
		return returnPages
	}
	return nil
}
func (p *Page) GetSubPages() []Page {
	if p.isType == DIRECTORY {
		var returnPages []Page
		for _, subPage := range p.SubPages {
			if subPage.isType == PAGE {
				returnPages = append(returnPages, subPage)
			}
		}
		return returnPages
	}
	return nil
}

// ReadDirectory recibe la ruta al directorio relativa a working directory y devuelve un Page
//completo listo para servir en un HandleFunc
func ReadDirectory(dirName string) (Page, error) {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return Page{}, err
	}

	dir := Page{
		WebPath:  "/" + dirName,
		filePath: dirName,
		isType:   DIRECTORY,
	}

	for _, file := range files {
		if file.IsDir() {
			subDir, errDir := ReadDirectory(filepath.Join(dir.filePath, file.Name()))
			if errDir != nil {
				return Page{}, errDir
			}
			dir.SubPages = append(dir.SubPages, subDir)
		} else {
			subPage, errPage := ReadPageFile(file, filepath.Join(dir.filePath, file.Name()))
			if errPage != nil {
				return Page{}, errPage
			}
			dir.SubPages = append(dir.SubPages, subPage)
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
		WebPath:  "/" + strings.TrimSuffix(dirName, filepath.Ext(dirName)),
		filePath: dirName,
		isType:   PAGE,
	}, nil
}

// HandleDirectory recorre recursivamente Page y todos los Page.
//SubPages ejecutando la función HandleFunc para servir los archivos
func HandleDirectory(dir *Page) error {
	if dir.isType != DIRECTORY {
		return errors.New("imposible manejar Page. No es un directorio")
	}

	// Con este bucle for, para recuperar la dirección de memoria correcta de SubPages[i],
	//con range se envía la misma dirección de memoria y todas las direcciones servirán el último
	//Page almacenado en esa dirección
	for i := 0; i < len(dir.SubPages); i++ {
		var err error
		switch dir.SubPages[i].isType {
		case PAGE:
			err = HandlePage(&dir.SubPages[i])
		case DIRECTORY:
			err = HandleDirectory(&dir.SubPages[i])

		}
		if err != nil {
			return err
		}
	}

	return nil
}

func HandlePage(p *Page) error {
	if p.isType != PAGE {
		return errors.New("imposible manejar Page. No es una página válida")
	}
	HandleFunc(
		p.WebPath, func(w http.ResponseWriter, r *http.Request) {
			if p.bufferedData == nil {
				data, err := ioutil.ReadFile(p.filePath)
				if err != nil {
					log.Fatalf("%s error leyendo archivo", p.filePath)
					return
				}
				if p.GetExtension() == ".md" {
					p.bufferedData = blackfriday.Run(data)
				} else {
					p.bufferedData = data
				}
			}
			w.Write(p.bufferedData)
		},
	)
	return nil
}
