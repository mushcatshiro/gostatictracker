package reminder

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TetsValidate(t *testing.T) {
	_, err := validate("*/", "minute", 0, 60)
	assert.Error(t, err)

	_, err = validate("*/1", "minute", 0, 60)
	assert.NoError(t, err)

	_, err = validate("0", "minute", 0, 60)
	assert.NoError(t, err)
	_, err = validate("1", "minute", 0, 60)
	assert.NoError(t, err)
	_, err = validate("59", "minute", 0, 60)
	assert.NoError(t, err)

	_, err = validate("60", "minute", 0, 60)
	assert.Error(t, err)
}
