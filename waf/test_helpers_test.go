package waf_test

import (
	"os"

	"github.com/spf13/afero"
)

func mockFs(c string) afero.Afero {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "main.tf", []byte(c), os.ModePerm)
	return afero.Afero{Fs: fs}
}
