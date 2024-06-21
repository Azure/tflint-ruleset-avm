// Package outputs provides the rules for the outputs category.
// Add the rules to the below slice to enable them.
package outputs

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var Rules = []tflint.Rule{
	// Removed as per spec change 20240621 https://github.com/Azure/Azure-Verified-Modules/pull/1028
	// NewRequiredOutputRule("required_output_tffr2", "resource", "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tffr2---category-outputs---additional-terraform-outputs"),
	NewRequiredOutputRule("required_output_rmfr7", "resource_id", "https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmfr7---category-outputs---minimum-required-outputs"),
}
