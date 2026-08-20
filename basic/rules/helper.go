package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func ForFiles(runner tflint.Runner, action func(tflint.Runner, *hcl.File) error) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := action(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}
