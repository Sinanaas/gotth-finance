package utils

import (
	"bytes"
	"html/template"
	"log"
)

const balanceTemplate = `<input
	type="text"
	value="Rp. {{.Balance}}"
	disabled
	aria-label="Account balance"
	class="rounded-lg mt-2 p-2 border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-amber-400 bg-white w-full"
/>`

func GetMessageTemplate(message string) []byte {
	tmpl, err := template.New("balance").Parse(balanceTemplate)
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
