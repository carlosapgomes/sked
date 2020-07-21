package web

import (
	"bytes"
	"html/template"

	"carlosapgomes.com/gobackend/internal/user"
)

type templateData struct {
	Title string
	User  *user.User
	Link  string
}

func (app App) render(fileName string, td *templateData) (*bytes.Buffer, error) {
	t, err := template.ParseFiles("./templates/" + fileName)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, td)
	if err != nil {
		return nil, err
	}
	// buf.WriteTo(wr)
	return buf, nil
}
