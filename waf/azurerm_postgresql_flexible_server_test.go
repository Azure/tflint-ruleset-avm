package waf_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermPostgreSqlFlexibleServerZoneRedundancy(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: waf.AzurermPostgreSqlFlexibleServerZoneRedundancy(),
			content: `
	variable "high_availability_mode" {
		type    = string
		default = "ZoneRedundant"
	}
	resource "azurerm_postgresql_flexible_server" "example" {
		high_availability {
			mode = var.high_availability_mode
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: waf.AzurermPostgreSqlFlexibleServerZoneRedundancy(),
			content: `
    variable "high_availability_mode" {
		type    = string
		default = "SameZone"
	}
	resource "azurerm_postgresql_flexible_server" "example" {
		high_availability {
			mode = var.high_availability_mode
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    waf.AzurermPostgreSqlFlexibleServerZoneRedundancy(),
					Message: "SameZone is an invalid attribute value of `mode` - expecting (one of) [ZoneRedundant]",
				},
			},
		},
		{
			name: "missing block",
			rule: waf.AzurermPostgreSqlFlexibleServerZoneRedundancy(),
			content: `
	resource "azurerm_postgresql_flexible_server" "example" {

	}`,
			expected: helper.Issues{
				{
					Rule:    waf.AzurermPostgreSqlFlexibleServerZoneRedundancy(),
					Message: "The attribute `mode` must be specified",
				},
			},
		},
	}

	filename := "main.tf"
	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{filename: tc.content})
			stub := gostub.Stub(&attrvalue.AppFs, mockFs(tc.content))
			defer stub.Reset()
			if err := tc.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}

func TestAzurermPostgreSqlFlexibleServerCustomMaintenanceSchedule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: waf.AzurermPostgreSqlFlexibleServerCustomMaintenanceSchedule(),
			content: `
	variable "maintenance_window" {
		type    = string
		default = "1"
	}
	resource "azurerm_postgresql_flexible_server" "example" {
		maintenance_window {
			day_of_week = var.maintenance_window
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: waf.AzurermPostgreSqlFlexibleServerCustomMaintenanceSchedule(),
			content: `
    variable "maintenance_window" {
		type    = string
		default = "20"
	}
	resource "azurerm_postgresql_flexible_server" "example" {
		maintenance_window {
			day_of_week = var.maintenance_window
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    waf.AzurermPostgreSqlFlexibleServerCustomMaintenanceSchedule(),
					Message: "20 is an invalid attribute value of `day_of_week` - expecting (one of) [0 1 2 3 4 5 6]",
				},
			},
		},
		{
			name: "missing block",
			rule: waf.AzurermPostgreSqlFlexibleServerCustomMaintenanceSchedule(),
			content: `
	resource "azurerm_postgresql_flexible_server" "example" {

	}`,
			expected: helper.Issues{
				{
					Rule:    waf.AzurermPostgreSqlFlexibleServerCustomMaintenanceSchedule(),
					Message: "The attribute `day_of_week` must be specified",
				},
			},
		},
	}

	filename := "main.tf"
	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{filename: tc.content})
			stub := gostub.Stub(&attrvalue.AppFs, mockFs(tc.content))
			defer stub.Reset()
			if err := tc.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}
