package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		want *Config
		kind interface{}
		name string
	}{
		{
			name: "returns *Config struct",
			kind: &Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New()
			require.NoError(t, err)

			if tt.kind != nil {
				assert.IsType(t, tt.kind, c)
			}
			if tt.want != nil {
				assert.Equal(t, tt.want, c)
			}
		})
	}
}
