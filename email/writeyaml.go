package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bjin01/jobmonitor/schedules"
	"gopkg.in/yaml.v3"
)

func WriteYaml(result *schedules.Jobstatus) (string, error) {

	if len(result.Completed) > 0 {

		filename := fmt.Sprintf("/srv/pillar/sumapatch/%s", result.YamlFileName)
		var list []string
		for _, l := range result.Completed {
			list = append(list, l.Hostname)
		}
		var yamldata = map[string][]string{result.YamlFileName: list}
		var b bytes.Buffer
		yamlEncoder := yaml.NewEncoder(&b)
		yamlEncoder.SetIndent(2)
		yamlEncoder.Encode(&yamldata)

		if _, err := os.Stat("/srv/pillar/sumapatch"); os.IsNotExist(err) {
			// path/to/whatever does not exist
			logger.Infof("Directory for yaml output file is missing: %s\n", "/srv/pillar/sumapatch")
			return "", err
		} else {

			err = ioutil.WriteFile(filename, b.Bytes(), 0644)
			if err != nil {
				logger.Infof("Unable to write data into the file: %v\n", filename)
				return "", err
			}
			logger.Infof("List of systems with completed jobs has been written to: %s\n", filename)

		}

		if len(result.Pending) > 0 {

			filename := fmt.Sprintf("/srv/pillar/sumapatch/%s", result.YamlFileName_Pending)
			var list []string
			for _, l := range result.Pending {
				list = append(list, l.Hostname)
			}
			var yamldata = map[string][]string{"pending_minions": list}
			var b bytes.Buffer
			yamlEncoder := yaml.NewEncoder(&b)
			yamlEncoder.SetIndent(2)
			yamlEncoder.Encode(&yamldata)

			if _, err := os.Stat("/srv/pillar/sumapatch"); os.IsNotExist(err) {
				// path/to/whatever does not exist
				logger.Infof("Directory for yaml output file is missing: %s\n", "/srv/pillar/sumapatch")
				return "", err
			} else {

				err = ioutil.WriteFile(filename, b.Bytes(), 0644)
				if err != nil {
					logger.Infof("Unable to write data into the file: %v\n", filename)
					return "", err
				}
				logger.Infof("List of systems with pending jobs has been written to: %s\n", filename)

			}
		}

		if len(result.Failed) > 0 {

			filename := fmt.Sprintf("/srv/pillar/sumapatch/%s", result.YamlFileName_Failed)
			var list []string
			for _, l := range result.Failed {
				list = append(list, l.Hostname)
			}
			var yamldata = map[string][]string{"failed_minions": list}
			var b bytes.Buffer
			yamlEncoder := yaml.NewEncoder(&b)
			yamlEncoder.SetIndent(2)
			yamlEncoder.Encode(&yamldata)

			if _, err := os.Stat("/srv/pillar/sumapatch"); os.IsNotExist(err) {
				// path/to/whatever does not exist
				logger.Infof("Directory for yaml output file is missing: %s\n", "/srv/pillar/sumapatch")
				return "", err
			} else {

				err = ioutil.WriteFile(filename, b.Bytes(), 0644)
				if err != nil {
					logger.Infof("Unable to write data into the file: %v\n", filename)
					return "", err
				}
				logger.Infof("List of systems with failed jobs has been written to: %s\n", filename)

			}
		}

		return filename, nil
	} else {
		return "", fmt.Errorf("No systems with completed jobs. Therefore no completed_ file written.")
	}

}
