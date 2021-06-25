package table

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColumnDefinition_Options(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		in       []ColumnOpt
		expected uint8
	}{
		{
			in:       []ColumnOpt{},
			expected: 0,
		},
		{
			in:       []ColumnOpt{IndexColumn(), HiddenColumn()},
			expected: 17,
		},
	}

	for _, tt := range tests {
		cd := TextColumn("foo", tt.in...)
		require.Equal(t, tt.expected, cd.Options())
	}
}
