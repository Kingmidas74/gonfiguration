package gonfiguration

import (
	"os"
	"testing"
)

func TestConfiguration_GetSectionWhenYamlFileOnly(t *testing.T) {
	builder := InitConfigurationBuilder()
	builder, err := builder.AddYamlFile("appsettings.yaml")
	if err != nil {
		t.Error(err)
	}
	config, err := builder.Build()
	if err != nil {
		t.Error(err)
	}

	section, err := config.GetSection("AnotherSettings")
	if err != nil {
		t.Error(err)
	}

	if section["Flag"] == nil {
		t.Error("Key Flag not found")
	}

	value := section["Flag"].(string)

	if value == "false" {
		t.Error("Flag is not correct")
	}
}

func TestConfiguration_GetSectionWhenEnvOnly(t *testing.T) {

	os.Setenv("AnotherSettings:Flag", "true")
	defer os.Unsetenv("AnotherSettings:Flag")

	builder := InitConfigurationBuilder()

	builder, err := builder.AddEnvironmentVariables()
	if err != nil {
		t.Error(err)
	}
	config, err := builder.Build()
	if err != nil {
		t.Error(err)
	}

	section, err := config.GetSection("AnotherSettings")
	if err != nil {
		t.Error(err)
	}

	if section["Flag"] == nil {
		t.Error("Key Flag not found")
	}

	value := section["Flag"].(string)

	if value == "false" {
		t.Error("Flag is not correct")
	}
}

func TestConfiguration_GetSectionWhenEnvOverrideYaml(t *testing.T) {

	os.Setenv("AnotherSettings:Flag", "false")
	defer os.Unsetenv("AnotherSettings:Flag")

	builder := InitConfigurationBuilder()

	builder, err := builder.AddYamlFile("appsettings.yaml")
	if err != nil {
		t.Error(err)
	}

	builder, err = builder.AddEnvironmentVariables()
	if err != nil {
		t.Error(err)
	}

	config, err := builder.Build()
	if err != nil {
		t.Error(err)
	}

	section, err := config.GetSection("AnotherSettings")
	if err != nil {
		t.Error(err)
	}

	if section["Flag"] == nil {
		t.Error("Key Flag not found")
	}

	value := section["Flag"].(string)

	if value == "true" {
		t.Error("Flag is not correct")
	}
}

func TestConfiguration_GetSectionWhenYamlOverrideEnv(t *testing.T) {

	os.Setenv("AnotherSettings:Flag", "false")
	defer os.Unsetenv("AnotherSettings:Flag")

	builder := InitConfigurationBuilder()

	builder, err := builder.AddEnvironmentVariables()
	if err != nil {
		t.Error(err)
	}

	builder, err = builder.AddYamlFile("appsettings.yaml")
	if err != nil {
		t.Error(err)
	}

	config, err := builder.Build()
	if err != nil {
		t.Error(err)
	}

	section, err := config.GetSection("AnotherSettings")
	if err != nil {
		t.Error(err)
	}

	if section["Flag"] == nil {
		t.Error("Key Flag not found")
	}

	value := section["Flag"].(string)

	if value == "false" {
		t.Error("Flag is not correct")
	}
}
