package swashengine

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"tsundere/packages/customization/swash"
	"tsundere/packages/customization/swash/evaluator"
	"tsundere/source/master/sessions"
)

var (
	Directory = "resources/branding/"
)

type SwashEngine struct {
	session *sessions.Session
}

func New(session *sessions.Session) *SwashEngine {
	return &SwashEngine{session: session}
}

func (s *SwashEngine) Execute(path string, newLine bool, elements map[string]any) error {
	tokenizer, err := swash.NewTokenizerSourcedFromFile(Directory + path)
	if err != nil {
		log.Printf("Unable to create tokenizer due to void within the file specified: %v", err)
	}

	if err := tokenizer.Parse(); err != nil {
		log.Printf("Unable to tokenize due to: %v", err)
	}

	var buffer = new(bytes.Buffer)

	eval := evaluator.NewEvaluator(tokenizer, buffer, s.session.Channel)

	for key, value := range elements {
		err := eval.Memory.Go2Swash(key, value)
		if err != nil {
			log.Printf("Failed to add element: %v", err)
			continue
		}
	}

	for s, m := range s.Packages() {
		err := eval.Memory.WritePackage(s, m)
		if err != nil {
			log.Printf("Failed to add package: %v", err)
			continue
		}
	}

	if err := eval.Execute(); err != nil {
		return s.session.Println(err.Error())
	}

	if len(buffer.String()) < 1 {
		return s.session.Print()
	}

	if newLine {
		return s.session.Print(buffer.String() + "\r\n")
	}

	return s.session.Print(buffer.String())
}

func (s *SwashEngine) ExecuteString(path string, elements map[string]any) string {
	tokenizer, err := swash.NewTokenizerSourcedFromFile(Directory + path)
	if err != nil {
		log.Printf("Unable to create tokenizer due to void within the file specified: %v", err)
	}

	if err := tokenizer.Parse(); err != nil {
		log.Printf("Unable to tokenize due to: %v", err)
	}

	var buffer = new(bytes.Buffer)

	eval := evaluator.NewEvaluator(tokenizer, buffer, s.session.Channel)

	for key, value := range elements {
		err := eval.Memory.Go2Swash(key, value)
		if err != nil {
			log.Printf("Failed to add element: %v", err)
			continue
		}
	}

	for s, m := range s.Packages() {
		err := eval.Memory.WritePackage(s, m)
		if err != nil {
			log.Printf("Failed to add package: %v", err)
			continue
		}
	}

	if err := eval.Execute(); err != nil {
		return err.Error()
	}

	return buffer.String()
}
