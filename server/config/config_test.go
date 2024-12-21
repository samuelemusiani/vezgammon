package config

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestParse(t *testing.T) {
	err := Parse("config.toml")
	assert.NilError(t, err)

	c := Get()

	assert.Equal(t, c.Server.Bind, ":8080")
	assert.Equal(t, c.Server.Domain, "localhost:8080")
	assert.Equal(t, c.Docker, true)

	Set(&Config{})
}

func TestErrorParse(t *testing.T) {
	err := Parse("not-exists.toml")
	assert.ErrorContains(t, err, "no such file or directory")

	err = Parse("config.go")
	assert.ErrorContains(t, err, "toml: expected character")
}
