package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	err := JSON("./data/json/a.json", "./data/json/b.json", "./data/json/c.json")
	assert.NoError(t, err)
}
