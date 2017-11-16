package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromHumanSize(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       string
		expected    uint64
		shouldError bool
	}{
		{
			desc:        "empty should err",
			input:       "",
			expected:    0,
			shouldError: true,
		},
		{
			desc:        "invalid should err",
			input:       "huahuh",
			expected:    0,
			shouldError: true,
		},
		{
			desc:        "1megabyte",
			input:       "1M",
			expected:    1024 * 1024,
			shouldError: false,
		},

		{
			desc:        "doesnt error with trailing and leading",
			input:       "  1M   \t\t\n",
			expected:    1024 * 1024,
			shouldError: false,
		},
	}

	var (
		err    error
		actual uint64
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err = FromHumanSize(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
