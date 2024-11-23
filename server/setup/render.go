package setup

import (
	"log"
	"text/template"
)

func RenderTemplate(tmplName string) (*template.Template, error) {
	// parsing requested template files
	template, err := template.ParseFiles(tmplName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// parsing layout files
	template.ParseFiles("./templates/auth.layout.tmpl")
	template.ParseFiles("./templates/base.layout.tmpl")

	return template, nil
}
