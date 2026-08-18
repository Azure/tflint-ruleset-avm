// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package modulecontent

import (
	"os"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/zclconf/go-cty/cty"
)

func mockFs(c string) afero.Afero {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "main.tf", []byte(c), os.ModePerm)
	return afero.Afero{Fs: fs}
}

func TestFetchBlocks(t *testing.T) {
	mockBlockFetcher := new(MockBlockFetcher)
	mockBlockFetcher.On("BlockType").Return("resource")
	mockBlockFetcher.On("LabelOne").Return("azapi_resource")
	mockBlockFetcher.On("LabelNames").Return([]string{"type", "name"})
	mockBlockFetcher.On("Attributes").Return([]string{"name", "type"})
	content := `
resource "azapi_resource" "test" {
	type = "testType@0000-00-00"
	name = "testName"
}`
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})
	stub := gostub.Stub(&AppFs, mockFs(content))
	defer stub.Reset()
	ctx, blocks, diags := FetchBlocks(mockBlockFetcher, runner)
	if diags.HasErrors() {
		t.Fatalf("FetchBlocks returned errors: %v", diags)
	}
	require.Len(t, blocks, 1)
	assert.Equal(t, "azapi_resource", blocks[0].Labels[0])
	ctx.EvaluateExpr(blocks[0].Body.Attributes["type"].Expr, cty.String) // nolint: errcheck
	val, diags := blocks[0].Body.Attributes["type"].Expr.Value(nil)
	assert.False(t, diags.HasErrors())
	assert.Equal(t, "testType@0000-00-00", val.AsString())
}

// MockBlockFetcher is a mock implementation of BlockFetcher for testing purposes.
type MockBlockFetcher struct {
	BlockFetcher
	mock.Mock
}

func (m *MockBlockFetcher) BlockType() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockBlockFetcher) LabelOne() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockBlockFetcher) LabelNames() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockBlockFetcher) Attributes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}
