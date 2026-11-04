package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Rules is a list of all rules
var Rules = []tflint.Rule{
	NewTerraformHeredocUsageRule(),
	NewTerraformModuleProviderDeclarationRule(),
	NewTerraformSensitiveVariableNoDefaultRule(),
	NewTerraformSensitiveFalseRule(),
	NewTerraformVariableNullableTrueRule(),
	NewTerraformModuleControlVariableRule(),
	NewTerraformLocalsFileRule(),
	NewTerraformDeprecatedVariableFileRule(),
	NewTerraformDeprecatedOutputFileRule(),
	NewTerraformVariableNegativeToggleRule(),
	NewTerraformRequiredVersionRule(),
	NewTerraformParentIDVariableRule(),
	NewAzapiPrimaryResourceNamingRule(),
	NewTerraformPatternModuleSourceRule(),
	NewTerraformResourceIDVariableValidationRule(),
	NewTerraformCollectionVariableNullableRule(),
	NewTerraformCountForEachUsageRule(),
}
