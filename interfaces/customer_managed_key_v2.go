package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var CustomerManagedKeyV2TypeString = `object({
	key_vault_key_uri = string
	user_assigned_identity = optional(object({
		client_id = string
	}), null)
})`

var customerManagedKeyV2Type = StringToTypeConstraintWithDefaults(CustomerManagedKeyV2TypeString)

var CustomerManagedKeyV2 = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(customerManagedKeyV2Type, cty.NullVal(cty.DynamicPseudoType), true),
	RuleName:      "customer_managed_key",
	VarTypeString: CustomerManagedKeyV2TypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#customer-managed-keys",
}
