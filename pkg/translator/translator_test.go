package translator_test

import (
	"testing"

	"github.com/famusovsky/md/pkg/translator"
	"github.com/stretchr/testify/assert"
)

func Test_Encrypt(t *testing.T) {
	t.Run("returns a string identifier", func(t *testing.T) {
		type testCase struct {
			id       int
			expected string
		}

		testCases := []testCase{
			{
				id:       1024,
				expected: "Mv",
			},
			{
				id:       0,
				expected: "y",
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
			assert.Equal(t, "Mv", res)
		}
	})
}

func Test_Translate(t *testing.T) {
	t.Run("returns an integer identifier", func(t *testing.T) {
		type testCase struct {
			str      string
			expected int
		}

		testCases := []testCase{
			{
				str:      "Mv",
				expected: 1024,
			},
			{
				str:      "y",
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
			res, err := translator.Translate("Mv")
			if err != nil {
				t.Error()
			}
			assert.Equal(t, 1024, res)
		}
	})
}
