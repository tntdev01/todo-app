package config

import (
	"fmt"
	"testing"
)

var testEnvironment = map[string]string{
	"APP_PORT":    "8081",
	"DB_HOST":     "db.example.test",
	"DB_PORT":     "5433",
	"DB_NAME":     "todos_test",
	"DB_USER":     "test_user",
	"DB_PASSWORD": "test_password",
	"DB_SSLMODE":  "require",
}

func setTestEnvironment(t *testing.T) {
	t.Helper()
	for key, value := range testEnvironment {
		t.Setenv(key, value)
	}
}

func TestLoad(t *testing.T) {
	setTestEnvironment(t)

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	want := &Config{
		AppPort:    "8081",
		DBHost:     "db.example.test",
		DBPort:     "5433",
		DBName:     "todos_test",
		DBUser:     "test_user",
		DBPassword: "test_password",
		DBSSLMode:  "require",
	}
	if *got != *want {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

func TestLoadPanicsWhenRequiredValueIsMissing(t *testing.T) {
	keys := []string{"APP_PORT", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSLMODE"}
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			setTestEnvironment(t)
			t.Setenv(key, "")

			defer func() {
				got := recover()
				if got != fmt.Sprintf("%s is required", key) {
					t.Errorf("Load() panic = %v, want %q", got, fmt.Sprintf("%s is required", key))
				}
			}()

			_, _ = Load()
		})
	}
}

func TestDatabaseURL(t *testing.T) {
	config := Config{
		DBHost:     "db.example.test",
		DBPort:     "5433",
		DBName:     "todos_test",
		DBUser:     "test_user",
		DBPassword: "test_password",
		DBSSLMode:  "require",
	}

	got := config.DatabaseURL()
	want := "postgres://test_user:test_password@db.example.test:5433/todos_test?sslmode=require"
	if got != want {
		t.Errorf("DatabaseURL() = %q, want %q", got, want)
	}
}
