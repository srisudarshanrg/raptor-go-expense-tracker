package models

// TemplateData is the struct containing the data types to be sent to the template
type TemplateData struct {
	Info  map[string]string
	Error map[string]string
	Data  map[string]interface{}
}
