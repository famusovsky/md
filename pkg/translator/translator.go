package translator

import (
	"errors"
	"strings"
)

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = len(alphabet)

func Translate(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty string")
	}

	var result int

	for _, c := range s {
		if !strings.Contains(alphabet, string(c)) {
			return 0, errors.New("wrong char within string: " + string(c))
		}

		result = result*alphabetLen + strings.IndexRune(alphabet, c)
	}

	return result, nil
}

func Encrypt(id int) string {
	var result string

	for id > 0 {
		result = string(alphabet[id%alphabetLen]) + result
		id = id / alphabetLen
	}

	// for len(result) < 5 {
	// 	result = string(alphabet[0]) + result
	// }

	return result
}
