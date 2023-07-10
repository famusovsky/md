package translator_test

import (
	"testing"

	"github.com/famusovsky/md/pkg/translator"
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
			actual := translator.Encrypt(tc.id)
			assert.Equal(t, tc.expected, actual)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			res := translator.Encrypt(1024)
			assert.Equal(t, "yyyMv", res)
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
			actual, err := translator.Translate(tc.str)
			if err != nil {
				t.Error()
			}
			assert.Equal(t, tc.expected, actual)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			res, err := translator.Translate("yyyMv")
			if err != nil {
				t.Error()
			}
			assert.Equal(t, 1024, res)
		}
	})
}
