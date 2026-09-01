// Package outputs provides the rules for the outputs category.
// Add the rules to the below slice to enable them.
package outputs

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var Rules = []tflint.Rule{
	NewRequiredOutputRule("avm_output_resource_id_required", "resource_id", "https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmfr7---category-outputs---minimum-required-outputs"),
	NewNoEntireResourceOutputRule("avm_output_entire_resource_disallowed", "https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-tffr2---category-outputs---additional-terraform-outputs"),
}
