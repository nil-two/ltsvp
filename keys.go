package main

func ParseKeysList(list string) []string {
	if list == "" || list == "\\" {
		return []string{}
	}

	var keys []string
	var buf []rune
	var isEscaping bool

	for _, ch := range list {
		if isEscaping {
			buf = append(buf, ch)
			isEscaping = false
			continue
		}

		if ch == '\\' {
			isEscaping = true
			continue
		}

		if ch == ',' {
			keys = append(keys, string(buf))
			buf = make([]rune, 0, 8)
			continue
		}

		buf = append(buf, ch)
	}

	keys = append(keys, string(buf))
	return keys
}
