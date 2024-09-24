package dockerfile

import (
	"fmt"
	"io"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func Parse(path string) (*Dockerfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%q: %v", path, err)
	}
	defer f.Close()
	return ParseReader(f)
}

func ParseReader(r io.Reader) (*Dockerfile, error) {
	result, err := parser.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("dockerfile/parser.Parse %v", err)
	}
	stages, metaArgs, err := instructions.Parse(result.AST, nil)
	if err != nil {
		return nil, fmt.Errorf("dockerfile/instructions.Parse %v", err)
	}
	var out Dockerfile
	for _, metaArg := range metaArgs {
		metaArgOut := &MetaArg{ArgCommand: metaArg}
		for _, kv := range metaArg.Args {
			metaArgOut.Key = kv.Key
			if defaultValue := kv.Value; defaultValue != nil {
				{
					defaultValueCopy := *defaultValue
					metaArgOut.DefaultValue = &defaultValueCopy
				}
				{
					defaultValueCopy := *defaultValue
					metaArgOut.Value = &defaultValueCopy
				}
			}
		}
		out.MetaArgs = append(out.MetaArgs, metaArgOut)
	}
	for _, stage := range stages {
		outStage := &Stage{Stage: stage}
		for _, command := range stage.Commands {
			outCommand := &Command{Command: command}
			outStage.Commands = append(outStage.Commands, outCommand)
		}
		out.Stages = append(out.Stages, outStage)
	}
	out.analyzeStages()
	return &out, nil
}
