package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
)

func TestLockTerraformVar(t *testing.T) {
	t.Parallel()
	fmt.Println(interfaces.DiagnosticSettings.TerrafromVar())
}
