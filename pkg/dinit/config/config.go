package config

import (
	"io/ioutil"

	"github.com/rosenlo/toolkits/log"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Services map[string]*Service
}
type Service struct {
	ExecStart   string            `yaml:"exec_start"`
	Lifecycle   *Lifecycle        `yaml:"lifecycle"`
	DependsOn   []string          `yaml:"depends_on"`
	Environment map[string]string `yaml:"environment"`
}

type Lifecycle struct {
	PreStart  *LifecycleHandler `yaml:"pre_start"`
	PostStart *LifecycleHandler `yaml:"post_start"`
	PreStop   *LifecycleHandler `yaml:"pre_stop"`
	PostStop  *LifecycleHandler `yaml:"post_stop"`
}

type LifecycleHandler struct {
	Exec struct {
		Command string `yaml:"command"`
	}
	Http struct {
		Host    string   `yaml:"host"`
		Path    string   `yaml:"path"`
		Headers []string `yaml:"headers"`
	}
}

func ParseConfig(fpath string) (config *Config, err error) {
	config = new(Config)
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Error(err)
		return
	}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
