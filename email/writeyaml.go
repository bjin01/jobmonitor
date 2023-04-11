package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
			log.Default().Printf("Directory for yaml output file is missing: %s\n", "/srv/pillar/sumapatch")
			return "", err
		} else {

			err = ioutil.WriteFile(filename, b.Bytes(), 0644)
			if err != nil {
				log.Default().Printf("Unable to write data into the file: %v\n", filename)
				return "", err
			}
			log.Default().Printf("List of systems with completed jobs has been written to: %s\n", filename)

		}
		return filename, nil
	} else {
		return "", fmt.Errorf("No systems with completed jobs. Therefore no completed_ file written.")
	}

}
