package main

import (
	"github.com/simulatedsimian/assert"
	"testing"
)

func TestGoGet(t *testing.T) {
	assert.NoError(t, assert.Pack(goget("github.com/simulatedsimian/gocmdutil", "tempgopath", "buildlog")))
}
