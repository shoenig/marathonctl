package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

type Format int

const (
	Human Format = iota
	Json
	JsonPP
	Raw
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
	case "raw":
		return Formatter{Raw}
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
	case Raw:
		return f.Raw(input)
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

// Raw returns the raw input unmodified.
func (f Formatter) Raw(input io.Reader) string {
	b, e := ioutil.ReadAll(input)
	Check(e == nil, "failed to read input", e)
	return string(b)
}

const maxCols = 100

// Columnize will pretty print columns of information
// rows are by newline
// columns are by whitespace
func Columnize(text string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	all := [][]string{}
	longests := [maxCols]int{} // index=col, val=maxlength

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for i, field := range fields {
			if len(field) > longests[i] {
				longests[i] = len(field)
			}
		}

		all = append(all, fields)
	}
	e := scanner.Err()
	Check(e == nil, "scanner error", e)
	return strings.TrimSpace(fmtFields(longests, all))
}

func fmtFields(longests [maxCols]int, matrix [][]string) string {
	var b bytes.Buffer
	for _, fields := range matrix {
		for col, field := range fields {
			b.WriteString(pad(longests[col], field))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func pad(length int, text string) string {
	var b bytes.Buffer
	b.WriteString(text)
	for i := 0; i < (2+length)-len(text); i++ {
		b.WriteString(" ")
	}
	return b.String()
}
