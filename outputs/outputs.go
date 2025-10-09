// Package outputs provides the rules for the outputs category.
// Add the rules to the below slice to enable them.
package outputs

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var Rules = []tflint.Rule{
	NewRequiredOutputRule("required_output_rmfr7", "resource_id", "https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmfr7---category-outputs---minimum-required-outputs"),
	NewNoEntireResourceOutputRule("no_entire_resource_output_tffr2", "https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-tffr2---category-outputs---additional-terraform-outputs"),
}
