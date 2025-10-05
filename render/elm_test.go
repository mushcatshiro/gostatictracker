package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttrsStructToString(t *testing.T) {
	// string
	a1 := attrsStruct{class: "test-class"}
	assert.Equal(t, `class="test-class"`, a1.toString())

	a2 := attrsStruct{class: "test-class", id: "test-id"}
	assert.Equal(t, `class="test-class" id="test-id"`, a2.toString())

	// bool
	a3 := attrsStruct{class: "test-class", required: true}
	assert.Equal(t, `class="test-class" required`, a3.toString())
}

