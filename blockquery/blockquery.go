// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blockquery

// BlockQuery is a struct that represents a query to run against a block with the given type and labels.
type BlockQuery struct {
	BlockType       string            // The type of block to query, e.g. "module", "resource", "data", etc.
	LabelOne        string            // The name of the block to query, e.g. "azapi_resource".
	BlockLabelNames []string          // The labels to use to identify the block, e.g. ["type", "name"].
	Query           string            // The gjson query to run against the block.
	QueryAttribute  string            // The attribute to query, e.g. "body".
	CompareFunc     ResultCompareFunc // The function to use to compare the result of the query.
}

// NewBlockQuery returns a new BlockQuery.
func NewBlockQuery(
	blockType, labelOne string,
	blockLabelNames []string,
	queryAttribute, query string,
	cmpFn ResultCompareFunc) BlockQuery {
	return BlockQuery{
		BlockType:       blockType,
		BlockLabelNames: blockLabelNames,
		LabelOne:        labelOne,
		QueryAttribute:  queryAttribute,
		Query:           query,
		CompareFunc:     cmpFn,
	}
}
