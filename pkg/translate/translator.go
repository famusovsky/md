package translate

import "strings"

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = len(alphabet)

func Translate(s string) int {
	var result int

	for _, c := range s {
		result = result*alphabetLen + strings.IndexRune(alphabet, c)
	}

	return result
}

func Encrypt(id int) string {
	var result string

	for id > 0 {
		result = string(alphabet[id%alphabetLen]) + result
		id = id / alphabetLen
	}

	for len(result) < 5 {
		result = string(alphabet[0]) + result
	}

	return result
}
