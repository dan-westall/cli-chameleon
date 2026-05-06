package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCreateTemplate(t *testing.T) {
	os.Remove(FileName)
	t.Cleanup(func() { os.Remove(FileName) })

	if err := CreateTemplate(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(FileName)
	if err != nil {
		t.Fatal(err)
	}

	// Should be valid YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("template is not valid YAML: %v", err)
	}
}
