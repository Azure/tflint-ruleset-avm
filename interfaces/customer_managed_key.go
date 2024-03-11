package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var CustomerManagedKeyTypeString = `object({
      key_vault_resource_id  = string
      key_name               = string
      key_version            = optional(string, null)
      user_assigned_identity = optional(object({
        resource_id = string
      }), null)
    })`

var customerManagedKeyType = StringToTypeConstraintWithDefaults(CustomerManagedKeyTypeString)

var CustomerManagedKey = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(customerManagedKeyType, cty.NullVal(cty.DynamicPseudoType), true),
	RuleName:      "customer_managed_key",
	VarTypeString: CustomerManagedKeyTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#customer-managed-keys",
}
