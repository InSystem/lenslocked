package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// NewView creates View
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View is type View
type View struct {
	Template *template.Template
	Layout   string
}

// ServeHttp is used to ... 
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render view with predefined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// Eg the input {"home"} would result in the output 
// {"views/home"} if TeplateDir == "views/"
func addTemplatePath(files [] string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// Eg the input {"home"} would result in the output 
// {"home.gohtml"} if TeplateExt == "gohtml"
func addTemplateExt(files [] string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
