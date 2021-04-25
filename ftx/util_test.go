package ftx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddOptions(t *testing.T) {
	type testOption struct {
		Foo string `url:"foo"`
	}

	const u = "https://example.com"
	tests := []struct {
		opts interface{}
		want string
	}{
		{
			opts: nil,
			want: "https://example.com",
		},
		{
			opts: testOption{Foo: "bar"},
			want: "https://example.com?foo=bar",
		},
	}
	for _, tt := range tests {
		got, err := addOptions(u, tt.opts)
		assert.NoError(t, err)
		assert.Equal(t, tt.want, got)
	}
}
