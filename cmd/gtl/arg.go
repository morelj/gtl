package main

import (
	"flag"
	"strings"
)

type multiStringValueFlag []string

func (m *multiStringValueFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiStringValueFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

var _ flag.Value = (*multiStringValueFlag)(nil)
