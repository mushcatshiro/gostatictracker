package server

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestCleanDuplicatedPath(t *testing.T) {
	s1 := "/editor/editor/fpath/fname"
	o1, ok := cleanDuplicatedPath(s1, "")
	if !ok {
		t.Fatalf("fail to clean duplicate on %s", s1)
	}
	assert.Equal(t, "/editor/fpath/fname", o1)

	s2 := "/editor/editor/editor/fpath/fname"
	o2, ok := cleanDuplicatedPath(s2, "")
	if !ok {
		t.Fatalf("fail to clean duplicate on %s", s2)
	}
	assert.Equal(t, "/editor/fpath/fname", o2)
}
