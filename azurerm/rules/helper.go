package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var headMetaArgPriority = map[string]int{"for_each": 0, "count": 0, "provider": 1}
var tailMetaArgPriority = map[string]int{"lifecycle": 0, "depends_on": 1}

// IsHeadMeta checks whether a name represents a type of head Meta arg
func IsHeadMeta(argName string) bool {
	_, isHeadMeta := headMetaArgPriority[argName]
	return isHeadMeta
}

// IsTailMeta checks whether a name represents a type of tail Meta arg
func IsTailMeta(argName string) bool {
	_, isTailMeta := tailMetaArgPriority[argName]
	return isTailMeta
}

// Check checks whether the tf config files match given rules
func Check(runner tflint.Runner, check func(tflint.Runner, *hcl.File) error) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := check(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}
