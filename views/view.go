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
	Data     *Data
	Layout   string
}

type Alert struct {
	Type    string
	Message string
}

type Data struct {
	Alert     Alert
	CSRFtoken string
	User      bool
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

	d := &Data{
		Alert:     Alert{Type: Success, Message: ""},
		CSRFtoken: "",
		User:      false,
	}

	return &View{
		Template: tmpl,
		Data:     d,
		Layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter) error {
	return v.Template.ExecuteTemplate(w, v.Layout, v.Data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := v.Render(w); err != nil {
		panic(err)
	}
}
