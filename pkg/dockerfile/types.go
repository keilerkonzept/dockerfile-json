package dockerfile

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

type Dockerfile struct {
	MetaArgs []*MetaArg
	Stages   []*Stage
}

type MetaArg struct {
	instructions.ArgCommand `json:"-"`
	Key                     string
	DefaultValue            *string `json:","`
	ProvidedValue           *string `json:","`
	Value                   *string `json:","`
}

type From struct {
	Stage   *FromStage `json:",omitempty"`
	Scratch bool       `json:",omitempty"`
	Image   *string    `json:",omitempty"`
}

type FromStage struct {
	Named *string `json:",omitempty"`
	Index int
}

type Command struct {
	instructions.Command
	Name string
}

type Stage struct {
	instructions.Stage
	Name     *string `json:"As,omitempty"`
	From     From
	Commands []*Command
}
