package encoders

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
)

type decoder func(v interface{}) error

type FileReader struct {
	Writer  io.Reader
	Decoder decoder
}

func NewDecoder(encoderID int, reader io.Reader) *FileReader {
	var e decoder

	switch encoderID {
	case Yaml:
		e = yaml.NewDecoder(reader).Decode
	case Json:
		e = json.NewDecoder(reader).Decode
	}

	return &FileReader{Writer: reader, Decoder: e}
}

func (e *FileReader) Decode(v interface{}) error {
	return e.Decoder(v)
}
