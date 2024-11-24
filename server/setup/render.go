package setup

import (
	"log"
	"net/http"
	"text/template"

	"github.com/srisudarshanrg/go-expense-tracker/server/models"
)

// RenderTemplate parses and executes a template passed to the function as a filename
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmplName string, templateData models.TemplateData) error {
	// parsing requested template files
	parsedTemplate, err := template.ParseFiles("./templates/" + tmplName)
	if err != nil {
		log.Println(err)
		return err
	}
	// parsing layout files
	parsedTemplate.ParseFiles("./templates/auth.layout.tmpl")
	parsedTemplate.ParseFiles("./templates/base.layout.tmpl")

	// executing template and passing template data
	err = parsedTemplate.Execute(w, templateData)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
