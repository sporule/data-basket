package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {

	assert.Equal(t, 2, min(3, 2), "It should return 2")
	assert.Equal(t, -1, min(-1, 0), "It should return -1")
}
