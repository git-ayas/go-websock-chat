package main

import (
	"bytes"
	"html/template"
)

type MesageDataType struct {
	Sender  string
	Content string
}

func GetParsedMessageComponent(data MesageDataType) (parsedTemplate []byte, err error) {

	var templateOutBuffer bytes.Buffer

	templateInst, err := template.ParseFiles("./views/components/message.html")
	if err != nil {
		return []byte{}, err
	}

	templateInst.Execute(&templateOutBuffer, data)

	return templateOutBuffer.Bytes(), err

}
