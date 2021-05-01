package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	var tm Time
	err := tm.UnmarshalJSON([]byte("1619858251.3361459"))

	assert.NoError(t, err)
	assert.Equal(t, "2021-05-01 16:37:31", tm.Time.Format("2006-01-02 15:04:05"))
}
