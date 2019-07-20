package main

import (
	"fmt"
	"strings"
)

// AssignmentsMap is a `flag.Value` for `KEY=VALUE` arguments.
type AssignmentsMap struct {
	Values map[string]*string
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *AssignmentsMap) Help() string {
	separator := "="
	return fmt.Sprintf("a key/value pair KEY[%sVALUE]", separator)
}

// Set is flag.Value.Set
func (fv *AssignmentsMap) Set(v string) error {
	separator := "="
	fv.Texts = append(fv.Texts, v)
	if fv.Values == nil {
		fv.Values = make(map[string]*string)
	}
	i := strings.Index(v, separator)
	if i < 0 {
		fv.Values[v[:i]] = nil
	}
	value := v[i+len(separator):]
	fv.Values[v[:i]] = &value
	return nil
}

func (fv *AssignmentsMap) String() string {
	return strings.Join(fv.Texts, ", ")
}
