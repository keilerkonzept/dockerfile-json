package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/keilerkonzept/dockerfile-json/pkg/dockerfile"
	"github.com/yalp/jsonpath"
)

var config struct {
	Quiet          bool
	Expand         bool
	JSONPathString string
	JSONPath       jsonpath.FilterFunc
	JSONPathRaw    bool
	BuildArgs      AssignmentsMap
	NonzeroExit    bool
}

var name = "dockerfile-json"
var version = "dev"
var jsonOut = json.NewEncoder(os.Stdout)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("[%s %s] ", filepath.Base(name), version))

	config.Expand = true
	flag.BoolVar(&config.Quiet, "quiet", config.Quiet, "suppress log output (stderr)")
	flag.BoolVar(&config.Expand, "expand-build-args", config.Expand, "expand build args")
	flag.StringVar(&config.JSONPathString, "jsonpath", config.JSONPathString, "select parts of the output using JSONPath (https://goessner.net/articles/JsonPath)")
	flag.BoolVar(&config.JSONPathRaw, "jsonpath-raw", config.JSONPathRaw, "when using JSONPath, output raw strings, not JSON values")
	flag.Var(&config.BuildArgs, "build-arg", config.BuildArgs.Help())
	flag.Parse()

	if config.Quiet {
		log.SetOutput(ioutil.Discard)
	}

	if flag.NArg() == 0 {
		flag.Usage()
	}

	if jsonPathString := config.JSONPathString; jsonPathString != "" {
		if jsonPathString[0] != '$' {
			jsonPathString = "$" + jsonPathString
		}
		jsonPath, err := jsonpath.Prepare(jsonPathString)
		if err != nil {
			log.Fatalf("parse jsonpath %s: %v", jsonPathString, err)
		}
		config.JSONPath = jsonPath
	}
}

func buildArgEnvExpander() dockerfile.SingleWordExpander {
	env := make(map[string]string, len(config.BuildArgs.Values))
	for key, value := range config.BuildArgs.Values {
		if value != nil {
			env[key] = *value
			continue
		}
		if value, ok := os.LookupEnv(key); ok {
			env[key] = value
		}
	}
	return func(word string) (string, error) {
		if value, ok := env[word]; ok {
			return value, nil
		}
		return "", fmt.Errorf("not defined: $%s", word)
	}
}

func main() {
	var dockerfiles []*dockerfile.Dockerfile
	for _, path := range flag.Args() {
		dockerfile, err := dockerfile.Parse(path)
		if err != nil {
			log.Printf("error: parse %q: %v", path, err)
			config.NonzeroExit = true
			continue
		}
		dockerfiles = append(dockerfiles, dockerfile)
	}
	if config.Expand {
		env := buildArgEnvExpander()
		for _, dockerfile := range dockerfiles {
			dockerfile.Expand(env)
		}
	}
	switch {
	case config.JSONPath != nil:
		for _, dockerfile := range dockerfiles {
			rawJSON, err := json.Marshal(dockerfile)
			if err != nil {
				log.Printf("error: evaluate jsonpath: %v", err)
				config.NonzeroExit = true
				continue
			}
			var data map[string]interface{}
			if err := json.Unmarshal(rawJSON, &data); err != nil {
				log.Printf("error: evaluate jsonpath: %v", err)
				config.NonzeroExit = true
				continue
			}
			result, err := config.JSONPath(data)
			if err != nil {
				log.Printf("error: evaluate jsonpath: %v", err)
				config.NonzeroExit = true
				continue
			}
			values, isArray := result.([]interface{})
			value, isString := result.(string)
			switch {
			case isString && config.JSONPathRaw:
				fmt.Println(value)
			case isArray && config.JSONPathRaw:
				for _, value := range values {
					fmt.Println(value)
				}
			case isArray && !config.JSONPathRaw:
				for _, value := range values {
					jsonOut.Encode(value)
				}
			default:
				jsonOut.Encode(result)
			}

		}
	default:
		for _, dockerfile := range dockerfiles {
			jsonOut.Encode(dockerfile)
		}
	}
	if config.NonzeroExit {
		os.Exit(1)
	}
}
