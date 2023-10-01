package spmigration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func writeMapToYAML(filename string, data map[string][]string) error {
	// Serialize the map to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	// Write the YAML data to the file
	err = ioutil.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return err
	}
	logger.Infoln("Wrote all minions in group to", filename)
	return nil
}
