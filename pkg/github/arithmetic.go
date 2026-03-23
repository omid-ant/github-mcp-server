package github

import (
	"context"
	"strconv"

	"github.com/github/github-mcp-server/pkg/inventory"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/github/github-mcp-server/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Calculate(t translations.TranslationHelperFunc) inventory.ServerTool {
	return NewTool(
		ToolsetMetadataArithmetic,
		mcp.Tool{
			Name:        "calculate",
			Description: t("TOOL_CALCULATE_DESCRIPTION", "Perform a basic arithmetic calculation. Supports addition, subtraction, multiplication, and division."),
			Annotations: &mcp.ToolAnnotations{
				Title:        t("TOOL_CALCULATE_TITLE", "Calculate"),
				ReadOnlyHint: true,
			},
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"operation": {
						Type:        "string",
						Description: t("TOOL_CALCULATE_OPERATION_DESCRIPTION", "The arithmetic operation to perform. One of: add, subtract, multiply, divide."),
						Enum:        []any{"add", "subtract", "multiply", "divide"},
					},
					"a": {
						Type:        "number",
						Description: t("TOOL_CALCULATE_A_DESCRIPTION", "The first operand."),
					},
					"b": {
						Type:        "number",
						Description: t("TOOL_CALCULATE_B_DESCRIPTION", "The second operand."),
					},
				},
				Required: []string{"operation", "a", "b"},
			},
		},
		nil,
		func(ctx context.Context, _ ToolDependencies, _ *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
			operation, err := RequiredParam[string](args, "operation")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			a, aOK, err := OptionalParamOK[float64](args, "a")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			if !aOK {
				return utils.NewToolResultError("missing required parameter: a"), nil, nil
			}
			b, bOK, err := OptionalParamOK[float64](args, "b")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			if !bOK {
				return utils.NewToolResultError("missing required parameter: b"), nil, nil
			}

			var result float64
			switch operation {
			case "add":
				result = a + b
			case "subtract":
				result = a - b
			case "multiply":
				result = a * b
			case "divide":
				if b == 0 {
					return utils.NewToolResultError("division by zero"), nil, nil
				}
				result = a / b
			default:
				return utils.NewToolResultError("unknown operation: " + operation), nil, nil
			}

			return utils.NewToolResultText(strconv.FormatFloat(result, 'g', -1, 64)), nil, nil
		},
	)
}
