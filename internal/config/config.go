package config

import (
	"io/fs"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/hadithopen-io/back/pkg/errors"
)

type Config struct {
	API         API         `yaml:"api"`
	DB          DB          `yaml:"db"`
	HTTP        HTTP        `yaml:"http"`
	MeiliSearch MeiliSearch `yaml:"meiliSearch"`
	OpenAI      OpenAI      `yaml:"openAI"`
}

type API struct {
	Host string `yaml:"host"`
}

type DB struct {
	Conn string `yaml:"conn"`
}

type HTTP struct {
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
}

type MeiliSearch struct {
	APIKey string `yaml:"masterKey"`
	Host   string `yaml:"host"`
}

type OpenAIKey string

func (OpenAIKey) Env() string { return "OPENAI_API_KEY" }

func (o OpenAIKey) String() string { return string(o) }

type OpenAI struct {
	Key OpenAIKey `yaml:"key"`
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
