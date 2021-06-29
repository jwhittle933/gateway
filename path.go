package gateway

import "strings"

type Path string

func NewPath(s string) Path {
	return Path(s).trimLeading()
}

func (p Path) Parts() []string {
	return strings.Split(string(p), "/")
}

func (p Path) trimLeading() Path {
	if len(p) > 0 && p[0] == '/' {
		return p[1:]
	}

	return p
}
