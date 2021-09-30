package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	viper.Set("startTag", "${{")
	viper.Set("endTag", "}}")
	err := mergeFile("test.yaml", "../data/lds1.yaml", "../data/lds2.yaml")
	assert.NoError(t, err)
}
