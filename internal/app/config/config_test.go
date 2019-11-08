package config

import "testing"

func TestLoadConfig(t *testing.T) {
	_, err := Load()
	if err != nil {
		t.Fatal(err)
	}
}
