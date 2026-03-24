package github

import (
	"context"
	"testing"

	"github.com/github/github-mcp-server/internal/toolsnaps"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Calculate(t *testing.T) {
	t.Parallel()

	serverTool := Calculate(translations.NullTranslationHelper)
	tool := serverTool.Tool
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	assert.Equal(t, "calculate", tool.Name)
	assert.True(t, tool.Annotations.ReadOnlyHint)

	tests := []struct {
		name               string
		requestArgs        map[string]any
		expectToolError    bool
		expectedToolErrMsg string
		expectedResult     string
	}{
		{
			name:           "add two numbers",
			requestArgs:    map[string]any{"operation": "add", "a": 1.0, "b": 1.0},
			expectedResult: "2",
		},
		{
			name:           "add with zero",
			requestArgs:    map[string]any{"operation": "add", "a": 0.0, "b": 5.0},
			expectedResult: "5",
		},
		{
			name:           "subtract",
			requestArgs:    map[string]any{"operation": "subtract", "a": 10.0, "b": 3.0},
			expectedResult: "7",
		},
		{
			name:           "multiply",
			requestArgs:    map[string]any{"operation": "multiply", "a": 4.0, "b": 3.0},
			expectedResult: "12",
		},
		{
			name:           "divide",
			requestArgs:    map[string]any{"operation": "divide", "a": 10.0, "b": 4.0},
			expectedResult: "2.5",
		},
		{
			name:           "negative result",
			requestArgs:    map[string]any{"operation": "subtract", "a": 3.0, "b": 10.0},
			expectedResult: "-7",
		},
		{
			name:               "divide by zero",
			requestArgs:        map[string]any{"operation": "divide", "a": 5.0, "b": 0.0},
			expectToolError:    true,
			expectedToolErrMsg: "division by zero",
		},
		{
			name:               "unknown operation",
			requestArgs:        map[string]any{"operation": "modulo", "a": 5.0, "b": 2.0},
			expectToolError:    true,
			expectedToolErrMsg: "unknown operation",
		},
		{
			name:               "missing operation",
			requestArgs:        map[string]any{"a": 1.0, "b": 2.0},
			expectToolError:    true,
			expectedToolErrMsg: "missing required parameter: operation",
		},
		{
			name:               "missing a",
			requestArgs:        map[string]any{"operation": "add", "b": 2.0},
			expectToolError:    true,
			expectedToolErrMsg: "missing required parameter: a",
		},
		{
			name:               "missing b",
			requestArgs:        map[string]any{"operation": "add", "a": 1.0},
			expectToolError:    true,
			expectedToolErrMsg: "missing required parameter: b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := stubDeps{}
			handler := serverTool.Handler(deps)

			request := createMCPRequest(tc.requestArgs)
			result, err := handler(ContextWithDeps(context.Background(), deps), &request)
			require.NoError(t, err)

			if tc.expectToolError {
				require.True(t, result.IsError, "expected tool call result to be an error")
				errorContent := getErrorResult(t, result)
				assert.Contains(t, errorContent.Text, tc.expectedToolErrMsg)
				return
			}

			require.False(t, result.IsError)
			textContent := getTextResult(t, result)
			assert.Equal(t, tc.expectedResult, textContent.Text)
		})
	}
}
