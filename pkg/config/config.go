package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"git.cryptic.systems/volker.raschek/civ/pkg/domain"
	"gopkg.in/yaml.v2"
)

type FileReader interface {
	ReadFile() (*domain.Config, error)
}

func NewFileReader(name string) FileReader {
	switch filepath.Ext(name) {
	case ".json":
		return &JSONReader{
			name: name,
		}
	default:
		return &YAMLReader{
			name: name,
		}
	}
}

type FileWriter interface {
	WriteFile(config *domain.Config) error
}

func NewFileWriter(name string) FileWriter {
	switch filepath.Ext(name) {
	case ".json":
		return &JSONWriter{
			name: name,
		}
	default:
		return &YAMLWriter{
			name: name,
		}
	}
}

type JSONReader struct {
	name string
}

func (jr *JSONReader) read(r io.Reader) (*domain.Config, error) {
	config := new(domain.Config)
	jsonDecoder := json.NewDecoder(r)

	err := jsonDecoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (jr *JSONReader) ReadFile() (*domain.Config, error) {
	f, err := os.Open(jr.name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return jr.read(f)
}

type JSONWriter struct {
	name string
}

func (jw *JSONWriter) WriteFile(config *domain.Config) error {
	f, err := os.Create(jw.name)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return jw.write(f, config)
}

func (jw *JSONWriter) write(w io.Writer, config *domain.Config) error {
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.SetIndent("", "  ")

	err := jsonEncoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
}

type YAMLReader struct {
	name string
}

func (yr *YAMLReader) ReadFile() (*domain.Config, error) {
	f, err := os.Open(yr.name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return yr.read(f)
}

func (yr *YAMLReader) read(r io.Reader) (*domain.Config, error) {
	config := new(domain.Config)
	yamlDecoder := yaml.NewDecoder(r)

	err := yamlDecoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

type YAMLWriter struct {
	name string
}

func (yw *YAMLWriter) WriteFile(config *domain.Config) error {
	f, err := os.Create(yw.name)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return yw.write(f, config)
}

func (yw *YAMLWriter) write(w io.Writer, config *domain.Config) error {
	err := yaml.NewEncoder(w).Encode(config)
	if err != nil {
		return err
	}
	return nil
}
