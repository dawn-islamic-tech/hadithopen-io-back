package config

import (
	"io/fs"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/hadithopen-io/back/pkg/errors"
)

type Config struct {
	DB          DB          `yaml:"db"`
	HTTP        HTTP        `yaml:"http"`
	MeiliSearch MeiliSearch `yaml:"meiliSearch"`
	OpenAI      OpenAI      `yaml:"openAI"`
	Auth        Auth        `yaml:"auth"`
	JWT         JWT         `yaml:"jwt"`
}

type DB struct {
	Conn string `yaml:"conn"`
}

type HTTP struct {
	Host              string        `yaml:"host"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
}

type MeiliSearch struct {
	APIKey string `yaml:"masterKey"`
	Host   string `yaml:"host"`
}

type Auth struct {
	Secret string `yaml:"secret"`
}

type JWT struct {
	ExpiresTime time.Duration `yaml:"expiresTime"`
	RefreshTime time.Duration `yaml:"refreshTime"`
	Secret      string        `yaml:"secret"`
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
