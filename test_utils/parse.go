package test_utils

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ReadMultiYaml(filename string) ([]map[string]interface{}, error) {
	expected, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ParseMultiYaml(expected)
}

func ReadMultiYamlIO(f io.ReadCloser) ([]map[string]interface{}, error) {
	expected, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ParseMultiYaml(expected)
}

func ParseMultiYaml(data []byte) ([]map[string]interface{}, error) {
	out := []map[string]interface{}{}

	decoder := yaml.NewDecoder(bytes.NewBuffer(data))
	for {
		var data map[string]interface{}
		if err := decoder.Decode(&data); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		out = append(out, data)
	}

	return out, nil
}
