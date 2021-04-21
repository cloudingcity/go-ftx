package ftx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAuth(t *testing.T) {
	const (
		key    = "api-key"
		secret = "api-secret"
	)
	c := New(WithAuth(key, secret))

	assert.Equal(t, key, c.key)
	assert.Equal(t, []byte(secret), c.secret)
}

func TestWithSubAccount(t *testing.T) {
	tests := []struct {
		account string
		want    string
	}{
		{account: "my-account", want: "my-account"},
		{account: "my/account", want: "my%2Faccount"},
	}
	for _, tt := range tests {
		c := New(WithSubAccount(tt.account))

		assert.Equal(t, tt.want, c.subAccount)
	}
}
