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

func Merge(values *map[string]interface{}, filenames ...string) error {
	var (
		startTag = viper.GetString("startTag")
		endTag   = viper.GetString("endTag")
	)
	for _, filename := range filenames {
		var marshaler Marshaler
		var override map[string]interface{}
		ext := filepath.Ext(filename)
		if len(ext) > 1 {
			ext = ext[1:]
		}
		switch ext {
		case "json":
			marshaler = &json{}
		case "yaml":
			marshaler = &yaml{}
		default:
			return fmt.Errorf("no support type: %s", ext)
		}
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
			return err
		}
		if err := marshaler.Unmarshal(bufferWrite.Bytes(), &override); err != nil {
			return err
		}
		if values == nil {
			values = &override
		} else {
			if err := mergo.Map(values, override, mergo.WithSliceDeepCopy, mergo.WithAppendSlice); err != nil {
				return err
			}
		}
	}
	return nil
}
