package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestMerge(t *testing.T) {
	viper.Set("startTag", "${{")
	viper.Set("endTag", "}}")
	mergeFile("test.yaml", "../data/cds1.yaml", "../data/cds2.yaml")
}
