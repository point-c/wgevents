package helpers

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

func UnmarshalYAML[T any](filename string) (v T, err error) {
	var f *os.File
	f, err = os.Open(filename)
	if err != nil {
		return v, err
	}
	defer func() { errors.Join(err, f.Close()) }()

	if err = yaml.NewDecoder(f).Decode(&v); err != nil {
		return v, err
	}
	return
}
