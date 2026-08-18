package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var Rules = []tflint.Rule{
	NewAzurermArgOrderRule(),
	NewAzurermResourceTagRule(),
}
