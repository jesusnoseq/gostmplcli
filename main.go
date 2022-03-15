package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

const (
	generalErrorCode     = 1
	misuseOfShellBuiltin = 2
)

type Config struct {
	render       string
	output       string
	templates    arrayFlags
	outputWriter io.Writer
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	conf, output, err := parseFlags(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(output)
		os.Exit(misuseOfShellBuiltin)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing arguments:", err)
		fmt.Fprint(os.Stderr, "output:", output)
		os.Exit(generalErrorCode)
	}

	if conf.output == "" {
		conf.outputWriter = os.Stdout
	} else {
		f, err := os.Create(output)
		defer f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating file:", err)
		}
		conf.outputWriter = f
	}

	err = generateTemplate(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error generating template:", err)
		os.Exit(misuseOfShellBuiltin)
	}
}

func parseFlags(progname string, args []string) (config Config, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var conf Config
	flags.StringVar(&conf.render, "r", "", "set the template name to render")
	flags.StringVar(&conf.output, "o", "", "set render file output or stdout by default")
	flags.Var(&conf.templates, "t", "set templates paths to load")

	err = flags.Parse(args)
	if err != nil {
		return Config{}, buf.String(), err
	}
	return conf, buf.String(), nil
}

func loadEnvData() map[string]string {
	envVars := os.Environ()
	result := map[string]string{}
	for _, env := range envVars {
		// env is
		envPair := strings.SplitN(env, "=", 2)
		result[envPair[0]] = envPair[1]
	}
	return result
}

func generateTemplate(conf Config) error {
	t := template.New("")
	for _, path := range conf.templates {
		_, err := t.ParseGlob(path)
		if err != nil {
			return fmt.Errorf("error loading templates from %s. %w", path, err)
		}
	}
	envData := loadEnvData()
	err := t.ExecuteTemplate(conf.outputWriter, conf.render, envData)
	if err != nil {
		return fmt.Errorf("error rendering template. %w", err)
	}
	return nil
}
