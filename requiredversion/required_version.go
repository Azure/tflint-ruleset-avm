// Package requiredversion provides the rules for the requiredversion category.
// Add the rules to the below slice to enable them.
package requiredversion

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var Rules = []tflint.Rule{
	NewRequiredVersionRule("required_version_tfnfr25", "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr25---category-code-style---verified-modules-requirements"),
}
