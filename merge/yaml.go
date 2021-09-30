package merge

import (
	"os"

	yamlV3 "gopkg.in/yaml.v3"
)

type yaml struct{}

func (y *yaml) Unmarshal(buf []byte, value interface{}) error {
	return yamlV3.Unmarshal(buf, value)
}

func (y *yaml) Marshal(value interface{}) ([]byte, error) {
	return yamlV3.Marshal(value)
}

func YAML(filenames ...string) error {
	readFileLength := len(filenames) - 1
	mergeFilename := filenames[readFileLength]
	y := &yaml{}
	v, err := merge(y, filenames[0:readFileLength]...)
	if err != nil {
		return err
	}
	buf, err := y.Marshal(&v)
	if err != nil {
		return err
	}
	f, err := os.Create(mergeFilename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(buf)
	return err
}

func YAMLValue(filenames ...string) (map[string]interface{}, error) {
	y := &yaml{}
	return merge(y, filenames...)
}
