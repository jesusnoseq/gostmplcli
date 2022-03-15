// Unit tests for flag parsing.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	var tests = []struct {
		args []string
		conf Config
	}{
		{[]string{},
			Config{},
		},
		{[]string{"duh"},
			Config{},
		},
		{[]string{"-t", "templatepath"},
			Config{templates: arrayFlags{"templatepath"}},
		},
		{[]string{"-t", "templatepathA", "-t", "templatepathB", "-t", "templatepathC"},
			Config{templates: arrayFlags{"templatepathA", "templatepathB", "templatepathC"}},
		},
		{[]string{"-r", "rendertemplate"},
			Config{render: "rendertemplate"},
		},
		{[]string{"-o", "outputfile"},
			Config{output: "outputfile"},
		},
		{[]string{"-t", "templatepathA", "-t", "templatepathB", "-r", "rendertemplate", "-o", "outputfile"},
			Config{output: "outputfile", render: "rendertemplate", templates: arrayFlags{"templatepathA", "templatepathB"}},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			conf, output, err := parseFlags("prog", tt.args)
			if err != nil {
				t.Errorf("err got %v, want nil", err)
			}
			if output != "" {
				t.Errorf("output got %q, want empty", output)
			}
			if !reflect.DeepEqual(conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v", conf, tt.conf)
			}
		})
	}
}

func TestParseFlagsUsage(t *testing.T) {
	var usageArgs = []string{"-help", "-h", "--help"}

	for _, arg := range usageArgs {
		t.Run(arg, func(t *testing.T) {
			conf, output, err := parseFlags("prog", []string{arg})
			if err != flag.ErrHelp {
				t.Errorf("err got %v, want ErrHelp", err)
			}
			if !reflect.ValueOf(conf).IsZero() {
				t.Errorf("conf got %v, want nil", conf)
			}
			if strings.Index(output, "Usage of") < 0 {
				t.Errorf("output can't find \"Usage of\": %q", output)
			}
		})
	}
}

func TestParseFlagsError(t *testing.T) {
	var tests = []struct {
		args   []string
		errstr string
	}{
		{[]string{"-foo"}, "flag provided but not defined"},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			conf, output, err := parseFlags("prog", tt.args)
			if !reflect.ValueOf(conf).IsZero() {
				t.Errorf("conf got %v, want nil", conf)
			}
			if strings.Index(err.Error(), tt.errstr) < 0 {
				t.Errorf("err got %q, want to find %q", err.Error(), tt.errstr)
			}
			if strings.Index(output, "Usage of prog") < 0 {
				t.Errorf("output got %q", output)
			}
		})
	}
}

func TestGenerateTemplateWithFiles(t *testing.T) {

	var testsCases = []struct {
		name   string
		conf   Config
		result string
		errStr string
	}{
		{"empty params", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{},
			render:       "",
		}, "", "error rendering template"},
		{"non existing templates", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{"no_template"},
			render:       "",
		}, "", "error loading template"},
		{"non existing render template", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{},
			render:       "no_render_template",
		}, "", "error rendering template"},
		{"non existing included template", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{"test_data/*.input"},
			render:       "template_fail.input",
		}, "Template not ok ", "error rendering template"},
		{"template with include", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{"test_data/*.input"},
			render:       "template_c.input",
		}, "a b c", ""},
		{"template with include overwrite", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{"test_data/*.input", "test_data/other_folder/*.input"},
			render:       "template_c.input",
		}, "other a template b c", ""},
	}

	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {
			err := generateTemplate(tt.conf)
			if err != nil && tt.errStr == "" {
				t.Errorf("got error %s when no error was expected", err)
			}
			if err == nil && tt.errStr != "" {
				t.Errorf("got no error but error %s was expected", tt.errStr)
			}
			if err != nil && tt.errStr != "" {
				if strings.Index(err.Error(), tt.errStr) < 0 {
					t.Errorf("got [%s] but [%s] was expected", err.Error(), tt.errStr)
				}
			}
			if tt.conf.outputWriter.(*bytes.Buffer).String() != tt.result {
				expected := tt.conf.outputWriter.(*bytes.Buffer).String()
				t.Errorf("got [%s] but [%s] was expected", expected, tt.result)
			}
		})
	}
}

func TestGenerateTemplateWithEnvVars(t *testing.T) {
	var testsCases = []struct {
		name   string
		conf   Config
		result string
	}{
		{"template with existing env var", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{"test_data/template_env_a.input"},
			render:       "template_env_a.input",
		}, "Template with an FOO var"},
		{"template with non existing env var", Config{
			outputWriter: &bytes.Buffer{},
			templates:    []string{"test_data/template_env_b.input"},
			render:       "template_env_b.input",
		}, "Template with an <no value> var"},
	}
	os.Setenv("bar", "FOO")
	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {
			err := generateTemplate(tt.conf)
			if err != nil {
				t.Errorf("Error %s not expected", err)
			}
			if tt.conf.outputWriter.(*bytes.Buffer).String() != tt.result {
				expected := tt.conf.outputWriter.(*bytes.Buffer).String()
				t.Errorf("got [%s] but [%s] was expected", expected, tt.result)
			}
		})
	}
	os.Unsetenv("bar")
}

func TestLoadEnvData(t *testing.T) {
	os.Setenv("FOO", "1")
	envVars := loadEnvData()
	if envVars["FOO"] != "1" {
		t.Errorf("FOO env var not found")
	}
	os.Unsetenv("FOO")
	envVars = loadEnvData()
	if _, exists := envVars["FOO"]; exists {
		t.Errorf("FOO env var not found")
	}
}
