package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/n-creativesystem/envoy-watch/logger"
	"github.com/n-creativesystem/envoy-watch/merge"
	"github.com/n-creativesystem/envoy-watch/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type setting struct {
	Output string   `yaml:"output"`
	Files  []string `yaml:"files"`
}

type settings struct {
	Settings []setting `yaml:"settings"`
}

func (s *settings) run() error {
	for _, setting := range s.Settings {
		if err := mergeFile(setting.Output, setting.Files...); err != nil {
			return err
		}
	}
	return nil
}

func watch(cmd *cobra.Command, args []string) error {
	errBuffer := &bytes.Buffer{}
	log := logrus.New()
	log.SetLevel(logger.GetLogLevel())
	log.SetOutput(errBuffer)
	log.SetReportCaller(true)

	flags := cmd.Flags()
	settingFile, _ := flags.GetString("setting")
	envoyConfig, _ := flags.GetString("envoy-config")
	envoyArgs, _ := flags.GetString("envoy-args")
	buf, err := os.ReadFile(settingFile)
	if err != nil {
		log.Errorln(err)
		return errors.New(errBuffer.String())
	}
	var s settings
	if err := yaml.Unmarshal(buf, &s); err != nil {
		log.Errorln(err)
		return errors.New(errBuffer.String())
	}

	fw := watcher.New()
	if err := s.run(); err != nil {
		log.Errorln(err)
		return errors.New(errBuffer.String())
	}
	for _, setting := range s.Settings {
		for _, f := range setting.Files {
			fw.Add(f, func(in fsnotify.Event) {
				logrus.Debugf("%s\n", in.String())
				if err := s.run(); err != nil {
					log.Errorln(err)
				}
			})
		}
	}

	execCmd, stdout, stderr, err := execEnvoy(envoyConfig, envoyArgs)
	if err != nil {
		log.Errorln(err)
		return errors.New(errBuffer.String())
	}
	cmdStdout := logger.NewStdLogger(cmd.OutOrStdout())
	cmdStderr := logger.NewStdLogger(cmd.OutOrStderr())

	go io.Copy(cmdStdout, stdout)
	go io.Copy(cmdStderr, stderr)
	if err := execCmd.Wait(); err != nil {
		log.Errorln(err)
		return errors.New(errBuffer.String())
	}
	return nil
}

func mergeFile(outputFile string, filenames ...string) error {
	value, err := merge.Merge(filenames...)
	if err != nil {
		return err
	}
	ext := filepath.Ext(outputFile)
	if len(ext) > 1 {
		ext = ext[1:]
	}
	var buf []byte
	switch ext {
	case "json":
		buf, err = json.Marshal(value)
	case "yaml":
		buf, err = yaml.Marshal(value)
	}
	if err != nil {
		return err
	}

	directory := filepath.Dir(outputFile)
	if _, err := os.Stat(directory); err != nil {
		_ = os.MkdirAll(directory, 0755)
	}
	tempName := fmt.Sprintf("%s.tmp", outputFile)
	logrus.Debugf("create temporary file: %s", tempName)
	if f, err := os.Create(tempName); err == nil {
		defer f.Close()
		_, err = f.Write(buf)
		if err != nil {
			return err
		}
		if err := os.Rename(f.Name(), outputFile); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func execEnvoy(config, args string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	logrus.Debug(fmt.Sprintf("envoy -c %s %s", config, args))
	cmd := exec.Command("sh", "-c", fmt.Sprintf("envoy -c %s %s", config, args))
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	return cmd, stdout, stderr, err
}

func NewCmdRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run application",
		Long:  "run application",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.SetLevel(logger.GetLogLevel())
			if err := watch(cmd, args); err != nil {
				cmd.PrintErr(err.Error())
			}
		},
	}
	flags := cmd.Flags()
	flags.StringP("setting", "s", "setting.yaml", "")
	flags.StringP("envoy-config", "c", "envoy.yaml", "envoy config")
	flags.StringP("envoy-args", "a", "", "envoy arguments")
	return cmd
}
