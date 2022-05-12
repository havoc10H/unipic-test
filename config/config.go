package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var (
	configPath = "config.yaml"

	Global  GlobalConfig
	Config  ServerConfig
	Default DefaultValues
)

type GlobalConfig struct {
	Config  ServerConfig  `yaml:"config"`
	Default DefaultValues `yaml:"default"`
}

type ServerConfig struct {
	NasaURL      string `yaml:"nasaurl"`
	ImageTemPath string `yaml:"imagetempath"`
	VideoTemPath string `yaml:"videotempath"`
	IndexPath    string `yaml:"indexpath"`
	ServerPort   string `yaml:"serverport"`
}

type DefaultValues struct {
	URL         string `yaml:"url"`
	Title       string `yaml:"title"`
	CopyRight   string `yaml:"copyright"`
	Explanation string `yaml:"explanation"`
}

func init() {
	Load(configPath)
}

// Load loads config information into Global
func Load(file string) (GlobalConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("failed to read config file: %v\n", err)
		return Global, err
	}

	err = yaml.Unmarshal(data, &Global)
	if err != nil {
		log.Printf("failed to parse config file: %v\n", err)
		return Global, err
	}

	Config = Global.Config
	Default = Global.Default

	return Global, nil
}
