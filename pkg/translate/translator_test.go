package translate_test

import (
	"testing"

	"github.com/famusovsky/md/pkg/translate"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	t.Run("returns an alphanumeric short identifier", func(t *testing.T) {
		type testCase struct {
			id       int
			expected string
		}

		testCases := []testCase{
			{
				id:       1024,
				expected: "yyyMv",
			},
			{
				id:       0,
				expected: "yyyyy",
			},
		}

		for _, tc := range testCases {
			actual := translate.Encrypt(tc.id)
			assert.Equal(t, tc.expected, actual)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			assert.Equal(t, "yyyMv", translate.Encrypt(1024))
		}
	})
}

func TestTranslate(t *testing.T) {
	t.Run("returns an integer identifier", func(t *testing.T) {
		type testCase struct {
			str      string
			expected int
		}

		testCases := []testCase{
			{
				str:      "yyyMv",
				expected: 1024,
			},
			{
				str:      "yyyyy",
				expected: 0,
			},
		}

		for _, tc := range testCases {
			actual := translate.Translate(tc.str)
			assert.Equal(t, tc.expected, actual)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			assert.Equal(t, 1024, translate.Translate("yyyMv"))
		}
	})
}
