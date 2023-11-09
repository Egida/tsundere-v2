package simpleconfig

import (
	"os"
	"tsundere/packages/utilities/simpleconfig/encoders"
)

func (s *SimpleConfig) encode(path string, overwrite bool, v interface{}) error {
	if !overwrite {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return err
		}
	}

	createdFile, err := os.Create(path)
	if err != nil {
		return err
	}

	err = encoders.NewEncoder(int(s.coder), createdFile).Encode(v)
	if err != nil {
		return err
	}

	return nil
}
