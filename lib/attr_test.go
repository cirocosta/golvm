package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLvAttr(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       string
		expected    *LvAttr
		shouldError bool
	}{
		{
			desc:        "fails with empty",
			input:       "",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "fails with not enough chars",
			input:       "aa",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "fails with more than enough chars",
			input:       "aaaaaaaaaaaasssssaaaaaa",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "fails with unexpected chars",
			input:       "Ç---------",
			expected:    &LvAttr{},
			shouldError: true,
		},
		{
			desc:        "fails with unexpected chars",
			input:       "a         ",
			expected:    &LvAttr{},
			shouldError: true,
		},
		{
			desc:  "proper volume type parsing",
			input: "t---------",
			expected: &LvAttr{
				VolumeType:                "thin pool",
				Permissions:               "-",
				AllocationPolicy:          "-",
				FixedMinor:                "-",
				State:                     "-",
				DeviceState:               "-",
				TargetType:                "-",
				OverrideNewBlocksWithZero: "-",
				VolumeHealth:              "-",
				SkipActivation:            "-",
			},
			shouldError: false,
		},
		{
			desc:  "proper overall parsing",
			input: "twi-aotz--",
			expected: &LvAttr{
				VolumeType:                "thin pool",
				Permissions:               "writeable",
				AllocationPolicy:          "inherited",
				FixedMinor:                "-",
				State:                     "active",
				DeviceState:               "open",
				TargetType:                "thin",
				OverrideNewBlocksWithZero: "overwrite by zero",
				VolumeHealth:              "-",
				SkipActivation:            "-",
			},
			shouldError: false,
		},
		{
			desc:  "proper thin parse",
			input: "Vwi-a-tz--",
			expected: &LvAttr{
				VolumeType:                "thin volume",
				Permissions:               "writeable",
				AllocationPolicy:          "inherited",
				FixedMinor:                "-",
				State:                     "active",
				DeviceState:               "-",
				TargetType:                "thin",
				OverrideNewBlocksWithZero: "overwrite by zero",
				VolumeHealth:              "-",
				SkipActivation:            "-",
			},
			shouldError: false,
		},
	}

	var (
		parsed *LvAttr
		err    error
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			parsed, err = ParseLvAttr(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, parsed)

			assert.Equal(t, tc.expected.VolumeType,
				parsed.VolumeType)
			assert.Equal(t, tc.expected.Permissions,
				parsed.Permissions)
			assert.Equal(t, tc.expected.AllocationPolicy,
				parsed.AllocationPolicy)
			assert.Equal(t, tc.expected.FixedMinor,
				parsed.FixedMinor)
			assert.Equal(t, tc.expected.State,
				parsed.State)
			assert.Equal(t, tc.expected.DeviceState,
				parsed.DeviceState)
			assert.Equal(t, tc.expected.TargetType,
				parsed.TargetType)
			assert.Equal(t, tc.expected.OverrideNewBlocksWithZero,
				parsed.OverrideNewBlocksWithZero)
			assert.Equal(t, tc.expected.VolumeHealth,
				parsed.VolumeHealth)
			assert.Equal(t, tc.expected.SkipActivation,
				parsed.SkipActivation)
		})
	}
}
