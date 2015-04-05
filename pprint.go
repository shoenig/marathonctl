package main

import (
	"bufio"
	"bytes"
	"strings"
)

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
