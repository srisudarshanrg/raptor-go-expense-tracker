package setup

import (
	"log"
	"text/template"
)

func RenderTemplate(tmplName string) (*template.Template, error) {
	// parsing requested template files
	parsedTemplate, err := template.ParseFiles(tmplName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// parsing layout files
	parsedTemplate.ParseFiles("./templates/auth.layout.tmpl")
	parsedTemplate.ParseFiles("./templates/base.layout.tmpl")

	return parsedTemplate, nil
}
