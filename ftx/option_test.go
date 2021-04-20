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

	assert.Equal(t, c.key, key)
	assert.Equal(t, c.secret, []byte(secret))
}

func TestWithSubAccount(t *testing.T) {
	const account = "my-account"

	c := New(WithSubAccount(account))

	assert.Equal(t, c.subAccount, account)
}
