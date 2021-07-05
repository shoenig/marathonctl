package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Properties map[string]string

func (p Properties) GetStringOr(s, alt string) string {
	v, exists := p[s]
	if !exists {
		return alt
	}
	return v
}

func ReadProperties(filename string) (Properties, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("unable to parse line: %q", line)
		}
		m[tokens[0]] = tokens[1]

	}
	if scErr := scanner.Err(); scErr != nil {
		return nil, scErr
	}

	return m, nil
}
