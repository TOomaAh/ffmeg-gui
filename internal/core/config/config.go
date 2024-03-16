package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AcceptedExtensions []string `yaml:"accepted_extensions"`
	FfmpegPath         string   `yaml:"ffmpeg_path"`
	FfprobePath        string   `yaml:"ffprobe_path"`
}

func NewConfig() *Config {
	configFile := "config.yaml"

	// if file does not exist, create it
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// write the default config
		config := Config{
			AcceptedExtensions: []string{".mp4", ".mkv", ".avi", ".mov", ".flv", ".wmv", ".webm"},
		}

		data, err := yaml.Marshal(&config)
		if err != nil {
			panic(err)
		}

		_, err = file.Write(data)
		if err != nil {
			panic(err)
		}

		return &config

	}

	// read the file
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	// unmarshal the file
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return &config

}

func (c *Config) IsAcceptedExtension(extension string) bool {
	for _, e := range c.AcceptedExtensions {
		if e == extension {
			return true
		}
	}
	return false
}

func (c *Config) Update() {
	configFile := "config.yaml"

	// read the file
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	// unmarshal the file
	err = yaml.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}

}
