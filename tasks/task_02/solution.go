package main

func reverseRunes(runes []rune) {
	i := 0
	j := len(runes) - 1

	for i < j {
		runes[i], runes[j] = runes[j], runes[i]
		i++
		j--
	}
}

func rotateRunes(s string, shift int) string {
	runes := []rune(s)
	n := len(runes)

	if n == 0 {
		return ""
	}

	shift %= n
	if shift < 0 {
		shift += n
	}

	reverseRunes(runes[:shift])
	reverseRunes(runes[shift:])
	reverseRunes(runes)

	return string(runes)
}
