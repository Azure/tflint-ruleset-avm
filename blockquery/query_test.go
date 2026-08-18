// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blockquery

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestQuery(t *testing.T) {
	ctyVal := cty.ObjectVal(map[string]cty.Value{
		"key": cty.StringVal("value"),
	})
	query := "key"
	result, err := Query(ctyVal, cty.DynamicPseudoType, query)
	require.NoError(t, err)
	require.Equal(t, "value", result.String())
}

func TestQueryCty(t *testing.T) {
	testCases := []struct {
		desc      string
		in        cty.Value
		query     string
		out       cty.Value
		expectErr bool
	}{
		{
			desc: "simple query",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.StringVal("value"),
			}),
			query:     "key",
			out:       cty.StringVal("value"),
			expectErr: false,
		},
		{
			desc: "two level query",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.ObjectVal(map[string]cty.Value{
					"key2": cty.StringVal("value"),
				}),
			}),
			query:     "key.key2",
			out:       cty.StringVal("value"),
			expectErr: false,
		},
		{
			desc: "non existent key",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.StringVal("value"),
			}),
			query:     "key.keyNotExist",
			out:       cty.StringVal("value"),
			expectErr: true,
		},
		{
			desc: "simple list",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.ListVal([]cty.Value{
					cty.NumberIntVal(1),
					cty.NumberIntVal(2),
				}),
			}),
			query: "key",
			out: cty.ListVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
			expectErr: false,
		},
		{
			desc: "simple list with hash wildcard",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.ListVal([]cty.Value{
					cty.NumberIntVal(1),
					cty.NumberIntVal(2),
				}),
			}),
			query: "key.#",
			out: cty.ListVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
			expectErr: false,
		},
		{
			desc: "complex list",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"key1": cty.NumberIntVal(1),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"key1": cty.NumberIntVal(1),
					}),
				}),
			}),
			query: "key.#.key1",
			out: cty.ListVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(1),
			}),
			expectErr: false,
		},
		{
			desc: "complex nested list",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"key1": cty.ListVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"key2": cty.NumberIntVal(1),
							}),
							cty.ObjectVal(map[string]cty.Value{
								"key2": cty.NumberIntVal(1),
							}),
						}),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"key1": cty.ListVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"key2": cty.NumberIntVal(1),
							}),
							cty.ObjectVal(map[string]cty.Value{
								"key2": cty.NumberIntVal(1),
							}),
						}),
					}),
				}),
			}),
			query: "key.#.key1.#.key2",
			out: cty.ListVal([]cty.Value{
				cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(1)}),
				cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(1)}),
			}),
			expectErr: false,
		},
		{
			desc: "unknown value",
			in: cty.ObjectVal(map[string]cty.Value{
				"key": cty.UnknownVal(cty.String),
			}),
			query:     "key",
			out:       cty.UnknownVal(cty.String),
			expectErr: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			out, err := QueryCty(tC.in, tC.query)
			if tC.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tC.out, out)
			}
		})
	}
}

func TestNextQuerySegment(t *testing.T) {
	testCases := []struct {
		desc      string
		in        string
		segment   string
		remaining string
	}{
		{
			desc:      "two part",
			in:        "a.b",
			segment:   "a",
			remaining: "b",
		},
		{
			desc:      "one part",
			in:        "a",
			segment:   "a",
			remaining: "",
		},
		{
			desc:      "three part",
			in:        "a.b.c",
			segment:   "a",
			remaining: "b.c",
		},
		{
			desc:      "empty",
			in:        "",
			segment:   "",
			remaining: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			segment, remaining := nextQuerySegment(tC.in)
			require.Equal(t, tC.segment, segment)
			require.Equal(t, tC.remaining, remaining)
		})
	}
}
