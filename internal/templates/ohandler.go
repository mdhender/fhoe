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

package templates

//// HandlerT represents a single template
//type HandlerT struct {
//	once sync.Once
//	name string
//	root string
//	layout string
//	yield string
//	overrides []string
//	tracer trace.Tracer
//	t *template.Template
//}
//
//func NewT(name, root, layout, yield string, tracer trace.Tracer, overrides ...string) http.HandlerFunc {
//	t := &HandlerT{name: name, yield: yield, root: root, layout: layout, tracer: tracer}
//	t.overrides = append(t.overrides, overrides...)
//	if err := t.Load(); err != nil {
//		panic(err)
//	}
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		t.tracer.Tracef("%s: parsing template\n", t.layout)
//		var ctx struct {
//			Site struct {
//				Title string
//			}
//		}
//		ctx.Site.Title = "joe"
//		if err := t.Load(); err != nil {
//			panic(err)
//		} else if err = t.t.Execute(w, &ctx); err != nil {
//			panic(err)
//		}
//		t.tracer.Tracef("%s: parsed template\n", t.layout)
//	}
//}
//
//func (t *HandlerT) Load() error {
//	fmt.Println("load --------------------------------------------------------")
//	var files []string
//
//	err := filepath.Walk(t.root, func(path string, info os.FileInfo, err error) error {
//		// avoid directories, files without the gohtml extension, and files in the yield path
//		if !(info.IsDir() || filepath.Ext(path) != ".gohtml" || strings.Contains(path, "yield")) {
//			files = append(files, path)
//		}
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//
//	// sort the list of templates files in an half-hearted attempt to ensure that we can reliably reproduce behavior.
//	// the driving example is two templates in the path that have the same name.
//	// we'd like to ensure that they're always loaded in the same order.
//	sort.Strings(files)
//
//	// append the list of listOfTemplateFiles that customize the layout for this view
//	files = append(files, t.overrides...)
//
//	// the optional yield file
//	if t.yield != "" {
//		files = append(files, filepath.Join(t.root, t.yield))
//	}
//
//	// the layout file
//	files = []string{filepath.Join(t.root, t.layout)}
//
//	t.t, err = template.ParseFiles(files...)
//	return err
//}
//
////func (t *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
////	t.once.Do(func() {t.t = template.Must(template.ParseFiles(t.filename))})
////}
