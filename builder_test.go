package main

import (
	//"fmt"
	"github.com/simulatedsimian/assert"
	"testing"
)

func TestPkgHash(t *testing.T) {
	h1 := generatePackageHash("github.com/user/name")
	h2 := generatePackageHash("github.com/user/name")
	h3 := generatePackageHash("github.com/user/namex")

	assert.Equal(t, h1, h2)
	assert.NotEqual(t, h2, h3)
}

func TestGoGet(t *testing.T) {
	var binfo BuildInfo

	assert.NoError(t, assert.Pack(goget("github.com/simulatedsimian/assert", "buildlog", &binfo)))
	assert.NoError(t, assert.Pack(gotest("github.com/simulatedsimian/assert", "testlog", &binfo)))
	assert.NoError(t, assert.Pack(gocover("github.com/simulatedsimian/assert", "coverlog", &binfo)))
}
