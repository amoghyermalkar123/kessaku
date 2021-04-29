package kessaku

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Api struct {
	Pool struct {
		Context string `yaml:"context"`
		Size    int    `yaml:"size"`
		Batch   struct {
			IsStarted  string `yaml:"isActive"`
			Timer      string `yaml:"timeout"`
			ManualStop string `yaml:"manual"`
		} `yaml:"batch"`
	} `yaml:"pool"`
}

func parseAPI() []OptionSetter {
	data, err := os.Open("../API.yaml")
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	defer data.Close()
	a := &Api{}

	d := yaml.NewDecoder(data)

	// Start YAML decoding from file
	if err := d.Decode(&a); err != nil {
		log.Panicf("decoding faled: %s", err)
	}

	options := make([]OptionSetter, 0)

	if a.Pool.Context == "yes" {
		options = append(options, WithContext())
	}
	options = append(options, WithPoolSize(int(a.Pool.Size)))
	options = append(options, WithBatch())

	return options
}
