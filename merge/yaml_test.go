package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAML(t *testing.T) {
	err := YAML("./data/yaml/a.yaml", "./data/yaml/b.yaml", "./data/yaml/c.yaml")
	assert.NoError(t, err)
}
