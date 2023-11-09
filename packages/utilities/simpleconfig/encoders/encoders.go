package encoders

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
)

type encoder func(v interface{}) error

type FileWriter struct {
	Writer  io.Writer
	Encoder encoder
}

func NewEncoder(encoderID int, writer io.Writer) *FileWriter {
	var e encoder

	switch encoderID {
	case Yaml:
		e = yaml.NewEncoder(writer).Encode
	case Json:
		jsonE := json.NewEncoder(writer)
		jsonE.SetEscapeHTML(false)
		jsonE.SetIndent("", "\t")

		e = jsonE.Encode
	}

	return &FileWriter{Writer: writer, Encoder: e}
}

func (e *FileWriter) Encode(v interface{}) error {
	return e.Encoder(v)
}
