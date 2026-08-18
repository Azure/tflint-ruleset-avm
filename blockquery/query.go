// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blockquery

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type QueryErrorNotFound struct {
	Query string
}

func (e *QueryErrorNotFound) Error() string {
	return fmt.Sprintf("attribute not found: %s", e.Query)
}

func (e *QueryErrorNotFound) Err() error {
	return fmt.Errorf("%w", e)
}

func NewQueryErrorNotFound(query string) *QueryErrorNotFound {
	return &QueryErrorNotFound{Query: query}
}

// Query takes a cty value and a gjson query string and returns the result.
// The cty.Value is marshalled to JSON and the query is run against the resulting JSON.
func Query(val cty.Value, ty cty.Type, query string) (gjson.Result, error) {
	jsonbytes, err := ctyjson.Marshal(val, ty)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("could not marshal cty value: %s", err)
	}
	return gjson.GetBytes(jsonbytes, "value."+query), nil
}

// QueryCty takes a cty value and a query string and returns the result.
// It supports basic query syntax for objects and lists.
// The query string is a dot-separated list of attribute names.
// The query string may contain a list index or the hash wildcard (#).
// The hash wildcard is used to query all elements of a list.
func QueryCty(val cty.Value, query string) (cty.Value, error) {
	segment, remaining := nextQuerySegment(query)
	if i, isList := querySegmentPertainsToList(segment); isList {
		return queryList(val, i, segment, remaining)
	}
	if ok := val.Type().IsObjectType() || val.Type().IsMapType(); !ok {
		return cty.NilVal, fmt.Errorf("query segments remain and value is not an object or map")
	}
	attrs := val.Type().AttributeTypes()
	if _, ok := attrs[segment]; !ok {
		return cty.NilVal, NewQueryErrorNotFound(query)
	}
	if remaining == "" {
		return val.GetAttr(segment), nil
	}
	return QueryCty(val.GetAttr(segment), remaining)
}

// queryList is a supporting function of QueryCty that handles list operations.
func queryList(val cty.Value, i int, segment, remaining string) (cty.Value, error) {
	if !val.Type().IsListType() && !val.Type().IsTupleType() {
		return cty.NilVal, fmt.Errorf("query segment %s is a list operation but value is not a list", segment)
	}
	// -1 means the query used the hash wildcard
	if i == -1 {
		result := make([]cty.Value, 0, val.LengthInt())
		it := val.ElementIterator()
		for it.Next() {
			_, v := it.Element()
			if remaining == "" {
				result = append(result, v)
				continue
			}
			q, err := QueryCty(v, remaining)
			if err != nil {
				return cty.NilVal, err
			}
			result = append(result, q)
		}
		return cty.ListVal(result), nil
	}
	if i >= val.LengthInt() {
		return cty.NilVal, fmt.Errorf("index %d out of bounds for list of length %d", i, val.LengthInt())
	}
	next := val.Index(cty.NumberIntVal(int64(i)))
	if remaining == "" {
		return next, nil
	}
	return QueryCty(next, remaining)
}

// nextQuerySegment splits a query string into the first segment and the remaining part.
func nextQuerySegment(query string) (string, string) {
	before, after, _ := strings.Cut(query, ".")
	return before, after
}

// querySegmentPertainsToList checks if a query segment is a list operation.
func querySegmentPertainsToList(segment string) (int, bool) {
	if segment == "#" {
		return -1, true
	}
	i, err := strconv.Atoi(segment)
	if err != nil {
		return 0, false
	}
	return i, true
}
