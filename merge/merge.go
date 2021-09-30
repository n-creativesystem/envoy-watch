package merge

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"github.com/valyala/fasttemplate"
)

var (
	jsonMarshal = &json{}
	yamlMarshal = &yaml{}
)

func GetMarshal(filename string) (Marshaler, error) {
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		ext = ext[1:]
	}
	switch ext {
	case "json":
		return jsonMarshal, nil
	case "yaml":
		return yamlMarshal, nil
	default:
		return nil, fmt.Errorf("no support file extension: %s", ext)
	}
}

type Marshaler interface {
	Unmarshal(buf []byte, value interface{}) error
	Marshal(value interface{}) ([]byte, error)
}

func merge(marshaler Marshaler, filenames ...string) (map[string]interface{}, error) {
	var (
		startTag = viper.GetString("startTag")
		endTag   = viper.GetString("endTag")
	)
	var resultValues map[string]interface{}
	for _, filename := range filenames {
		var override map[string]interface{}
		bs, err := os.ReadFile(filename)
		if err != nil {
			continue
		}
		bufferWrite := &bytes.Buffer{}
		tpl := fasttemplate.New(string(bs), startTag, endTag)
		_, err = tpl.ExecuteFunc(bufferWrite, func(w io.Writer, tag string) (int, error) {
			v := os.Getenv(tag)
			return w.Write([]byte(v))
		})
		if err != nil {
			return nil, err
		}
		if err := marshaler.Unmarshal(bufferWrite.Bytes(), &override); err != nil {
			return nil, err
		}
		if resultValues == nil {
			resultValues = override
		} else {
			if err := mergo.Map(&resultValues, override, mergo.WithSliceDeepCopy, mergo.WithAppendSlice); err != nil {
				return nil, err
			}
		}
	}
	return resultValues, nil
}

func Merge(filenames ...string) (map[string]interface{}, error) {
	var (
		startTag = viper.GetString("startTag")
		endTag   = viper.GetString("endTag")
		valueTyp = map[string]interface{}{}
	)
	for _, filename := range filenames {
		marshaler, err := GetMarshal(filename)
		if err != nil {
			return nil, err
		}
		var override map[string]interface{}
		bs, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		bufferWrite := &bytes.Buffer{}
		tpl := fasttemplate.New(string(bs), startTag, endTag)
		_, err = tpl.ExecuteFunc(bufferWrite, func(w io.Writer, tag string) (int, error) {
			v := os.Getenv(tag)
			return w.Write([]byte(v))
		})
		if err != nil {
			return nil, err
		}
		if err := marshaler.Unmarshal(bufferWrite.Bytes(), &override); err != nil {
			return nil, err
		}
		if err := mergo.Map(&valueTyp, override, mergo.WithSliceDeepCopy, mergo.WithAppendSlice); err != nil {
			return nil, err
		}
	}
	delete(valueTyp, "aliases")
	return valueTyp, nil
}
