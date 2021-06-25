package table

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable_Spec(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name     string
		columns  []ColumnDefinition
		expected string
	}{
		{
			name:    "simple",
			columns: []ColumnDefinition{TextColumn("simple_text")},
			expected: `
{
  "name": "simple",
  "cacheable": false,
  "evented": false,
  "columns":[
    { "name": "simple_text", "type": "TEXT", "index": false, "required": false, "additional": false, "optimized": false, "hidden": false }
  ]
}`,
		},
	}

	mockGenerate := func(_ context.Context, _ QueryContext) ([]map[string]string, error) { return nil, nil }

	for _, tt := range tests {
		testTable := NewPlugin(tt.name, tt.columns, mockGenerate)
		generatedSpec, err := testTable.Spec()
		require.NoError(t, err, "generating spec for %s", tt.name)
		helperJSONEqVal(t, tt.expected, generatedSpec, "spec for %s", tt.name)
	}
}

func helperJSONEqVal(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) {
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		require.Fail(t, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()), msgAndArgs...)
		return
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		require.Fail(t, fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()), msgAndArgs...)
		return
	}

	require.EqualValues(t, expectedJSONAsInterface, actualJSONAsInterface, msgAndArgs...)
	return
}
