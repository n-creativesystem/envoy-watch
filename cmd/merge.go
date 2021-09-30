package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/n-creativesystem/envoy-watch/merge"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewCmdMerge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "merge file",
		Long:  "merge file",
		Run: func(cmd *cobra.Command, args []string) {
			filenames, _ := cmd.Flags().GetStringArray("files")
			outputFile, _ := cmd.Flags().GetString("output")
			jsonFiles := []string{}
			yamlFiles := []string{}
			for _, file := range filenames {
				ext := filepath.Ext(file)
				switch ext {
				case ".json":
					jsonFiles = append(jsonFiles, file)
				case ".yaml":
					yamlFiles = append(yamlFiles, file)
				}
			}
			var err error
			var yamlValue, jsonValue map[string]interface{}
			if len(yamlFiles) > 0 {
				yamlValue, err = merge.YAMLValue(yamlFiles...)
				if err != nil {
					logrus.Fatalln(err)
				}
			}
			if len(jsonFiles) > 0 {
				jsonValue, err = merge.JSONValue(jsonFiles...)
			}
			if err := mergo.Map(&yamlValue, jsonValue, mergo.WithAppendSlice, mergo.WithSliceDeepCopy); err != nil {
				logrus.Fatalln(err)
			}
			ext := filepath.Ext(outputFile)
			var buf []byte
			switch ext {
			case ".json":
				buf, err = json.Marshal(yamlValue)
			case ".yaml":
				buf, err = yaml.Marshal(yamlValue)
			}
			if err != nil {
				logrus.Fatalln(err)
			}
			if f, err := os.Create(outputFile); err == nil {
				defer f.Close()
				_, err = f.Write(buf)
				if err != nil {
					logrus.Fatalln(err)
				}
			} else {
				logrus.Fatalln(err)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringArrayP("files", "f", []string{}, "merge files")
	flags.StringP("output", "o", "", "output file name")
	return cmd
}
