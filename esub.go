package esub

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Eval(tmpl string, envList []string) (string, error) {
	envs := map[string]string{}

	for _, e := range envList {
		kv := strings.SplitN(e, "=", 2)
		envs[kv[0]] = kv[1]
	}

	reader := bufio.NewReader(strings.NewReader(tmpl))
	var buf strings.Builder

L:
	for {
		b, _ := reader.ReadBytes('$')

		if len(b) == 0 {
			break L
		}

		buf.Write(b[0 : len(b)-1])

		if b[len(b)-1] != '$' {
			buf.WriteByte(b[len(b)-1])
			break L
		}

		b, _ = reader.Peek(2)

		switch len(b) {
		case 0:
			buf.WriteByte('$')
			break L
		case 1:
			if b[0] == '{' {
				return "", fmt.Errorf("syntax error: %s", tmpl)
			} else {
				buf.WriteByte('$')
				continue
			}
		case 2:
			if b[0] == '$' && b[1] == '{' {
				buf.Write([]byte{'$', '{'})
				reader.Discard(2) //nolint:errcheck
				continue
			} else if b[0] != '{' {
				buf.WriteByte('$')
				continue
			}
		}

		reader.Discard(1) //nolint:errcheck
		b, _ = reader.ReadBytes('}')

		if len(b) == 0 {
			return "", fmt.Errorf("syntax error: %s", tmpl)
		}

		envName := string(b[0 : len(b)-1])

		if len(envName) == 0 || b[len(b)-1] != '}' {
			return "", fmt.Errorf("syntax error: %s", tmpl)
		}

		e, ok := envs[envName]

		if !ok {
			return "", fmt.Errorf("env '%s' not found: %s", envName, tmpl)
		}

		buf.WriteString(e)
	}

	return buf.String(), nil
}

func Fill(tmpl string) (string, error) {
	return Eval(tmpl, os.Environ())
}
