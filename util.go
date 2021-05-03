package main

import (
	"strings"
)

type Strings []string

func (ss *Strings) UnmarshalParam(s string) error {
	*ss = strings.Split(s, ",")
	return nil
}

func NormalizeSymbol(symbol string) string {
	return strings.ToLower(strings.TrimSpace(symbol))
}
