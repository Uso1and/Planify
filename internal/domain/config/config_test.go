package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

)

func TestGetConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		config   *ConfigDB
		expected string
	}{
		{
			name: "full connection string",
			config: &ConfigDB{
				Host:     "localhost",
				Port:     "5432",
				User:     "user",
				Password: "pass",
				Name:     "dbname",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=user password=pass dbname=dbname sslmode=disable",
		},
		{
			name: "empty values",
			config: &ConfigDB{
				Host:     "",
				Port:     "",
				User:     "",
				Password: "",
				Name:     "",
				SSLMode:  "",
			},
			expected: "host= port= user= password= dbname= sslmode=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetConnectionString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnv(t *testing.T) {

	originalValue := os.Getenv("TEST_ENV_VAR")
	defer os.Setenv("TEST_ENV_VAR", originalValue)

	t.Run("env variable set", func(t *testing.T) {
		os.Setenv("TEST_ENV_VAR", "test_value")
		result := getEnv("TEST_ENV_VAR", "default")
		assert.Equal(t, "test_value", result)
	})

	t.Run("env variable not set", func(t *testing.T) {
		os.Unsetenv("TEST_ENV_VAR")
		result := getEnv("TEST_ENV_VAR", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("env variable empty", func(t *testing.T) {
		os.Setenv("TEST_ENV_VAR", "")
		result := getEnv("TEST_ENV_VAR", "default")
		assert.Equal(t, "default", result)
	})
}
