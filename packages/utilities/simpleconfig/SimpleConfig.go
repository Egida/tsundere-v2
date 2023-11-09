package simpleconfig

import "tsundere/packages/utilities/simpleconfig/encoders"

type Coder int

const (
	Json = encoders.Json
	Yaml = encoders.Yaml
)

type SimpleConfig struct {
	directory     string
	coder         Coder
	fileExtension string
}

func New(coder Coder, directory string) *SimpleConfig {
	var fileExtension string
	switch coder {
	case Json:
		fileExtension = ".json"
	case Yaml:
		fileExtension = ".yml"
	}

	return &SimpleConfig{directory: directory, coder: coder, fileExtension: fileExtension}
}
