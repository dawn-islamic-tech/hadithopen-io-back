package config

import (
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/hadithopen-io/back/pkg/errors"
)

type Config struct {
	API API
}

type API struct {
	Host string `yaml:"host"`
}

func NewConfig(path string) (_ *Config, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "open config by path %s", path)
	}
	defer errors.AppendCloser(&err, file)

	var config Config
	return &config, errors.Wrap(
		yaml.NewDecoder(file).Decode(&config),
		"decode config information",
	)
}
