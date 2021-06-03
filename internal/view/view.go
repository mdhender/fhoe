// fhoe - Far Horizons order entry
// Copyright (c) 2021 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package view implements the view for the application.
// It implements a handler for templates.
// A view is the highest level handler. It is responsible for loading all the
// templates used in the view. It also handles the http requests (which it should not).
package view

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// View implements a view.
// It stores the templates to use for rendering.
type View struct {
	Name  string
	Yield string // this is the optional override for content

	// path and list both hold templates files to load.
	// we load from the path first, then from the list.
	// that allows templates to be overridden.
	// when two templates have the same name, the template library discards the first one loaded, keeping only the last template with that name.
	layoutFile          string
	pathToTemplateFiles string
	listOfTemplateFiles []string

	// templateSet holds the set of parsed templates.
	templateSet *template.Template
}

// New creates a new view.
// It will load all the template listOfTemplateFiles along with an optional list of listOfTemplateFiles.
// The intent of that list is to "specialize" the view by providing templates that have a common name but different view logic.
// I think that we pretty much avoid using them.
func New(name string, layoutFile string, yield string, pathToTemplateFiles string, listOfSpecializedTemplates ...string) (*View, error) {
	log.Printf("[views] creating new view: %q\n", name)
	v := &View{
		Name:                name,
		Yield:               yield,
		layoutFile:          layoutFile,
		pathToTemplateFiles: pathToTemplateFiles,
	}
	v.listOfTemplateFiles = append(v.listOfTemplateFiles, listOfSpecializedTemplates...)
	return v, nil
}

// Load loads all template files for the view.
// It starts by loading all files found in pathToTemplateFiles.
// It then loads the optional listOfTemplateFiles to customize.
//
// In a production environment, we would cache the loads.
// For now, though, we reload every time.
// This allows us to easily test changes to templates.
func (v *View) Load() (*template.Template, error) {
	fmt.Println(v)
	return v.reload()
}

// Render takes a context (the data parameter) and uses that to render the resulting view.
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	t, err := v.Load()
	if err != nil {
		fmt.Println(v)
		fmt.Printf("[%s] load: %+v\n", v.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if err = t.ExecuteTemplate(w, v.Name, data); err != nil {
		fmt.Printf("[%s] exec: %+v\n", v.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// reload loads all template files for the view.
// It starts by loading all files found in pathToTemplateFiles.
// It then loads the optional listOfTemplateFiles to customize.
func (v *View) reload() (*template.Template, error) {
	fmt.Println("reload --------------------------------------------------------")
	var files []string
	err := filepath.Walk(v.pathToTemplateFiles, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".gohtml" {
			// avoid files in the yield path
			if !strings.Contains(path, "yield") {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// sort the list of templates files in an half-hearted attempt to ensure that we can reliably reproduce behavior.
	// the driving example is two templates in the path that have the same name.
	// we'd like to ensure that they're always loaded in the same order.
	sort.Strings(files)

	// append the list of listOfTemplateFiles that customize the layout for this view
	files = append(files, v.listOfTemplateFiles...)

	// the optional yield file
	if v.Yield != "" {
		files = append(files, v.pathToTemplateFiles+v.Yield)
	}

	// the layout file
	files = append(files, v.layoutFile)

	return template.ParseFiles(files...)
}
