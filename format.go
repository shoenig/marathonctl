package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

type Format int

const (
	Human Format = iota
	Json
	JsonPP
)

type Humanize func(input io.Reader) string

type Formatter struct {
	format Format
}

func NewFormatter(f string) Formatter {
	switch f {
	case "jsonpp":
		return Formatter{JsonPP}
	case "json":
		return Formatter{Json}
	default:
		return Formatter{Human}
	}
}

func (f Formatter) Format(input io.Reader, h Humanize) string {
	switch f.format {
	case JsonPP:
		return f.JsonPPize(input)
	case Json:
		return f.Jsonize(input)
	default:
		return h(input)
	}
}

// Jsonize returns raw response on one line with no extra space.
func (f Formatter) Jsonize(input io.Reader) string {
	var s bytes.Buffer
	b, e := ioutil.ReadAll(input)
	Check(e == nil, "failed to read input", e)
	json.Compact(&s, b)
	return s.String()
}

// JsonPPize takes the raw response and adds newlines and indentations.
func (f Formatter) JsonPPize(input io.Reader) string {
	var s bytes.Buffer
	b, e := ioutil.ReadAll(input)
	Check(e == nil, "failed to read input", e)
	json.Indent(&s, b, "", "    ")
	return s.String()
}
