package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Parse parses a configuration file into an object; returns an error if necessary.
func Parse(filename string, obj interface{}) error {
	configFile, err := os.Open(filename) // #nosec
	if err != nil {
		return err
	}
	defer configFile.Close()

	return yaml.NewDecoder(configFile).Decode(obj)
}
