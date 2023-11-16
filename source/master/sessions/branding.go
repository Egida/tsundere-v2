package sessions

import (
	"bytes"
	"fmt"
	"path/filepath"
	"tsundere/packages/customization/gradient"
	"tsundere/packages/customization/swash"
	"tsundere/packages/customization/swash/evaluator"
)

// ExecuteBranding will execute a swash script
func (s *Session) ExecuteBranding(objects map[string]any, branding ...string) error {
	// create new tokenizer
	tokenizer, err := swash.NewTokenizerSourcedFromFile(filepath.Join(append([]string{"resources", "branding"}, branding...)...))
	if err != nil {
		return err
	}

	// parse tokens
	if err := tokenizer.Parse(); err != nil {
		return err
	}

	var buffer = new(bytes.Buffer)

	// create new evaluator >.<
	eval := evaluator.NewEvaluator(tokenizer, buffer, s.Terminal.Channel)
	if err := s.sync(eval); err != nil {
		return err
	}

	// iterates through all specific custom objects
	for key, value := range objects {
		err := eval.Memory.Go2Swash(key, value)
		if err == nil {
			continue
		}

		return err
	}

	// finally execute it
	if err := eval.Execute(); err != nil {
		return err
	}

	if buffer.Len() < 1 {
		return s.Print(buffer)
	}

	return s.Println(buffer)
}

// ExecuteBrandingToString will execute a swash script into a string
func (s *Session) ExecuteBrandingToString(objects map[string]any, branding ...string) (string, error) {
	var buffer = new(bytes.Buffer)

	// create new tokenizer
	tokenizer, err := swash.NewTokenizerSourcedFromFile(filepath.Join(append([]string{"resources", "branding"}, branding...)...))
	if err != nil {
		return "", err
	}

	// parse tokens
	if err := tokenizer.Parse(); err != nil {
		return "", err
	}

	// create new evaluator >.<
	eval := evaluator.NewEvaluator(tokenizer, buffer, buffer)
	if err := s.sync(eval); err != nil {
		fmt.Println(err)
		return "", err
	}

	// iterates through all specific custom objects
	for key, value := range objects {
		err := eval.Memory.Go2Swash(key, value)
		if err != nil {
			fmt.Println(err)
			return "", err
		}

	}

	// finally execute it
	if err := eval.Execute(); err != nil {
		fmt.Println(err)
		return "", err
	}

	return buffer.String(), nil
}

// ExecuteBrandingToStringNoError will execute a swash script into a string with no error
func (s *Session) ExecuteBrandingToStringNoError(objects map[string]any, branding ...string) string {
	literal, err := s.ExecuteBrandingToString(objects, branding...)
	if err != nil {
		return err.Error()
	}

	return literal
}

func (s *Session) sync(evaluator *evaluator.Evaluator) error {
	// built-in swash functions
	var builtIn = map[string]any{
		"user":     s.UserProfile,
		"gradient": gradient.New,

		"sessions": map[string]any{
			"length": Count(),
		},
	}

	// optional packages for swash
	var packages = map[string]map[string]any{
		"test": {
			"test": "this is a test. lol",
		},
	}

	// add all std thingies 2 swash
	for key, value := range builtIn {
		if err := evaluator.Memory.Go2Swash(key, value); err != nil {
			return err
		}
	}

	// add all packages 2 swash
	for key, value := range packages {
		if err := evaluator.Memory.WritePackage(key, value); err != nil {
			return err
		}
	}

	return nil
}
