package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Rules is a list of all rules
var Rules = []tflint.Rule{
	NewTerraformCountIndexUsageRule(),
	NewTerraformHeredocUsageRule(),
	NewTerraformLocalsOrderRule(),
	NewTerraformModuleProviderDeclarationRule(),
	NewTerraformOutputOrderRule(),
	NewTerraformOutputSeparateRule(),
	NewTerraformRequiredProvidersDeclarationRule(),
	NewTerraformRequiredVersionDeclarationRule(),
	NewTerraformResourceDataArgLayoutRule(),
	NewTerraformSensitiveVariableNoDefaultRule(),
	NewTerraformVariableNullableFalseRule(),
	NewTerraformVariableOrderRule(),
	NewTerraformVariableSeparateRule(),
	NewTerraformVersionsFileRule(),
}
