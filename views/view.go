package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	Error   string = "error"
	Success string = "success"
)

type View struct {
	Template *template.Template
	Layout   string
}

type Alert struct {
	Type    string
	Message string
}

func layoutFiles() []string {
	files, err := filepath.Glob("views/layouts/" + "*" + ".html")
	if err != nil {
		panic(err)
	}
	return files
}

func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: tmpl,
		Layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}
