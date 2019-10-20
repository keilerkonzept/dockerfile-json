package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type dockerfile struct {
	MetaArgs []instructions.ArgCommand
	Stages   []instructions.Stage
}

var config struct {
	Quiet     bool
	Expand    bool
	BuildArgs AssignmentsMap
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
	flag.Var(&config.BuildArgs, "build-arg", config.BuildArgs.Help())
	flag.Parse()

	if config.Quiet {
		jsonOut = json.NewEncoder(ioutil.Discard)
	}
}

func buildArgEnvExpander() instructions.SingleWordExpander {
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
	var dockerfiles []dockerfile
	for _, path := range flag.Args() {
		func() {
			f, err := os.Open(path)
			defer f.Close()
			if err != nil {
				log.Printf("error: %q: %v", path, err)
				return
			}
			result, err := parser.Parse(f)
			if err != nil {
				log.Printf("error: parse %q: %v", path, err)
				return
			}
			stages, metaArgs, err := instructions.Parse(result.AST)
			if err != nil {
				log.Printf("error: parse %q: %v", path, err)
				return
			}
			dockerfile := dockerfile{
				MetaArgs: metaArgs,
				Stages:   stages,
			}
			dockerfiles = append(dockerfiles, dockerfile)
		}()
	}
	if config.Expand {
		for _, dockerfile := range dockerfiles {
			dockerfile.expand(buildArgEnvExpander())
		}
	}

	type outCommand struct {
		Name string
		instructions.Command
	}

	type outStage struct {
		Name       string
		BaseName   string
		SourceCode string
		Platform   string
		IsDerived  bool
		IsScratch  bool
		Commands   []outCommand
	}

	type outDockerfile struct {
		MetaArgs []instructions.ArgCommand
		Stages   []outStage
	}

	for _, dockerfile := range dockerfiles {
		var out outDockerfile
		out.MetaArgs = dockerfile.MetaArgs
		seenStageNames := make(map[string]bool)
		for _, stage := range dockerfile.Stages {
			outStage := outStage{
				Name:       stage.Name,
				BaseName:   stage.BaseName,
				SourceCode: stage.SourceCode,
				Platform:   stage.Platform,
			}
			switch {
			case seenStageNames[stage.BaseName]:
				outStage.IsDerived = true
			case stage.BaseName == "scratch":
				outStage.IsScratch = true
			}
			if stage.Name != "" {
				seenStageNames[stage.Name] = true
			}
			for _, command := range stage.Commands {
				outStage.Commands = append(outStage.Commands, outCommand{
					Name:    command.Name(),
					Command: command,
				})
			}
			out.Stages = append(out.Stages, outStage)
		}
		jsonOut.Encode(out)
	}
}

func (d *dockerfile) expand(envExpander instructions.SingleWordExpander) {
	metaArgsEnvExpander := d.metaArgsEnvExpander(envExpander)
	for i, stage := range d.Stages {
		d.Stages[i].BaseName = os.Expand(stage.BaseName, func(key string) string {
			value, err := metaArgsEnvExpander(key)
			if err != nil {
				return ""
			}
			return value
		})
		for i := range stage.Commands {
			cmdExpander, ok := stage.Commands[i].(instructions.SupportsSingleWordExpansion)
			if ok {
				cmdExpander.Expand(metaArgsEnvExpander)
			}
		}
	}
}

func (d *dockerfile) metaArgsEnvExpander(envExpander instructions.SingleWordExpander) instructions.SingleWordExpander {
	metaArgsEnv := make(map[string]string, len(d.MetaArgs))
	for _, arg := range d.MetaArgs {
		if arg.Value != nil {
			metaArgsEnv[arg.Key] = *arg.Value
		}
		if value, err := envExpander(arg.Key); err == nil {
			arg.Value = &value
			metaArgsEnv[arg.Key] = value
		}
		err := arg.Expand(envExpander)
		if err != nil {
			continue
		}
		if arg.Value != nil {
			metaArgsEnv[arg.Key] = *arg.Value
		}
	}
	return func(key string) (string, error) {
		if value, ok := metaArgsEnv[key]; ok {
			return value, nil
		}
		return "", fmt.Errorf("not defined: $%s", key)
	}
}
