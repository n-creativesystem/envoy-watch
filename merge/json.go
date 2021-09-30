package merge

import (
	encjson "encoding/json"
	"os"
)

type json struct{}

func (j *json) Unmarshal(buf []byte, value interface{}) error {
	return encjson.Unmarshal(buf, value)
}

func (j *json) Marshal(value interface{}) ([]byte, error) {
	return encjson.Marshal(value)
}

func JSON(filenames ...string) error {
	readFileLength := len(filenames) - 1
	mergeFilename := filenames[readFileLength]
	j := &json{}
	v, err := merge(j, filenames[0:readFileLength]...)
	if err != nil {
		return err
	}
	buf, err := j.Marshal(&v)
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

func JSONValue(filenames ...string) (map[string]interface{}, error) {
	j := &json{}
	return merge(j, filenames...)
}
