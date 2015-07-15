package xml

import (
	"strconv"
)

func line2map(line string) (m map[string]string) {
	m = make(map[string]string)
	line_length := len(line)
	for i := 0; i < line_length; i++ {
		var key, val string
		if line[i] == ' ' || line[i] == '\t' {
			continue
		}
		for j := i; j < line_length; j++ {
			if line[j] == '=' {
				key = line[i:j]
				i = j + 1
				break
			}
		}
		if key == "" {
			break
		}

		sep := line[i]
		for j := i + 1; j < line_length; j++ {
			if line[j] == sep {
				val = line[i+1 : j]
				i = j + 1
				break
			}
		}
		m[key] = val
	}
	return
}

func str2int64(s string) int64 {
	if s == "" {
		return 0
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func str2float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func str2bool(s string) bool {
	if s == "" {
		s = "true" // visible='true'
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return b
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
