package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Config
		kind interface{}
	}{
		{
			name: "returns *Config struct",
			kind: &Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()

			if tt.kind != nil {
				assert.IsType(t, tt.kind, c)
			}
			if tt.want != nil {
				assert.Equal(t, tt.want, c)
			}
		})
	}
}
