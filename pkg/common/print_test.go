package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorIndex(t *testing.T) {

	pc := NewPrintColor()
	for i := 0; i < 6; i++ {
		pc.Get()("ola")
	}

	assert.Equal(t, pc.ColorIdx, 0)
}
