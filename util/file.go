package util

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml3 "gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type TagSliceCmdYamlMap []map[string][]string

func (p *TagSliceCmdYamlMap) UnmarshalYAML(value *yaml3.Node) error {
	if value.Kind != yaml3.MappingNode {
		return fmt.Errorf("pipeline must contain YAML mapping, has %v", value.Kind)
	}

	for i := 0; i < len(value.Content); i += 2 {
		var a string
		var b []string
		if err := value.Content[i].Decode(&a); err != nil {
			return err
		}
		if err := value.Content[i+1].Decode(&b); err != nil {
			return err
		}
		*p = append(*p, map[string][]string{a: b})
	}
	return nil
}

func TagGetYaml3(filepath string) ([]map[string][]string, error) {

	if !IsFileExist(filepath) {
		return nil, errors.Wrap(errors.New(filepath+" does not exist"), "errCode_1070")
	}
	if CheckEmptyFile(filepath) == false {
		f, err := ioutil.ReadFile(filepath)
		if err != nil {
			return nil, errors.Wrap(err, "can not find aws credentials at "+filepath+"\nPlease make sure aws client is installed and profile is setup appropriately.  See doc at...")
		}
		var y TagSliceCmdYamlMap
		if err = yaml3.Unmarshal(f, &y); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", y)
		return y, nil

	} else {
		log.Info(filepath + " is empty.  Nothing to process. Exiting...")

	}
	return nil, nil
}

func TagRemoveFile(filepath string) {
	if IsFileExist(filepath) {
		err := os.Remove(filepath)
		Check(err)
	}
}
