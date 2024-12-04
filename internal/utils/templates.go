package utils

import (
	"bytes"
	"html/template"
	"log"
)

func GetMessageTemplate(message string) []byte {
	tmpl, err := template.ParseFiles("internal/templates/balance.html")
	if err != nil {
		log.Println("template parsing error:", err)
		return nil
	}

	data := struct {
		Balance string
	}{
		Balance: message,
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, data)
	if err != nil {
		log.Println("template execution error:", err)
		return nil
	}

	return renderedMessage.Bytes()
}
