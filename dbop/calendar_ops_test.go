package dbop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCalendarRenderRange(t *testing.T) {
	_, err := getCalendarRenderRange(13, 2006)
	assert.EqualError(t, err, "Month must not be greater than 12")

	_, err = getCalendarRenderRange(0, -1)
	assert.EqualError(t, err, "Year must be greater or equal to 0")

	_, err = getCalendarRenderRange(1, 0)
	assert.EqualError(t, err, "Does not support unspecified year with specific month given")

	_, err = getCalendarRenderRange(0, 1)
	assert.NoError(t, err)
	_, err = getCalendarRenderRange(1, 1)
	assert.NoError(t, err)
	_, err = getCalendarRenderRange(0, 0)
	assert.NoError(t, err)
}
