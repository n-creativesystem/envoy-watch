package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "watch config file",
		Long:  "watch config file",
	}
	cobra.OnInitialize(initConfig)
	flags := cmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.watch.yaml)")
	flags.String("startTag", "${{", "environment variable start tag")
	flags.String("endTag", "}}", "environment variable end tag")
	flags.String("logLevel", "debug", "output log level")
	_ = viper.BindPFlag("startTag", flags.Lookup("startTag"))
	_ = viper.BindPFlag("endTag", flags.Lookup("endTag"))
	_ = viper.BindPFlag("logLevel", flags.Lookup("logLevel"))
	cmd.AddCommand(NewCmdRun())
	cmd.AddCommand(NewCmdMerge())
	cmd.AddCommand(NewCmdVersion())
	return cmd
}

func Execute() {
	cmd := NewCmdRoot()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".watch")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
