package models

// TemplateData is the struct containing the data types to be sent to the template
type TemplateData struct {
	Info     interface{}
	Error    interface{}
	Data     interface{}
	PostData interface{}
}
