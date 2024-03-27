package attrvalue

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

func toPtr[T any](val T) any {
	if reflect.TypeOf(val).Kind() == reflect.Pointer {
		return val
	}
	return &val
}

func toCtyType(val any) (cty.Type, error) {
	switch val.(type) {
	case *string:
		return cty.String, nil
	case string:
		return cty.String, nil
	case *int:
		return cty.Number, nil
	case int:
		return cty.Number, nil
	case float64:
		return cty.Number, nil
	case *float64:
		return cty.Number, nil
	case bool:
		return cty.Bool, nil
	case *bool:
		return cty.Bool, nil
	case []int32:
		return cty.List(cty.Number), nil
	case *[]int32:
		return cty.List(cty.Number), nil
	case []int:
		return cty.List(cty.Number), nil
	case *[]int:
		return cty.List(cty.Number), nil
	case []int64:
		return cty.List(cty.Number), nil
	case *[]int64:
		return cty.List(cty.Number), nil
	case []float32:
		return cty.List(cty.Number), nil
	case *[]float32:
		return cty.List(cty.Number), nil
	case []float64:
		return cty.List(cty.Number), nil
	case *[]float64:
		return cty.List(cty.Number), nil
	}
	return cty.NilType, fmt.Errorf("unsupported type %s", reflect.TypeOf(val).String())
}

func toCtyValue(val any) (cty.Value, error) {
	switch val.(type) {
	case string:
		s := val.(string)
		return cty.StringVal(s), nil
	case int:
		i := val.(int)
		return cty.NumberIntVal(int64(i)), nil
	case float64:
		f := val.(float64)
		return cty.NumberVal(big.NewFloat(f)), nil
	case bool:
		b := val.(bool)
		return cty.BoolVal(b), nil
	case []int32:
		ints := val.([]int32)
		var ctyInts []cty.Value
		for _, i := range ints {
			ctyInts = append(ctyInts, cty.NumberIntVal(int64(i)))
		}
		return cty.ListVal(ctyInts), nil
	case []int:
		ints := val.([]int)
		var ctyInts []cty.Value
		for _, i := range ints {
			ctyInts = append(ctyInts, cty.NumberIntVal(int64(i)))
		}
		return cty.ListVal(ctyInts), nil
	case []int64:
		ints := val.([]int64)
		var ctyInts []cty.Value
		for _, i := range ints {
			ctyInts = append(ctyInts, cty.NumberIntVal(i))
		}
		return cty.ListVal(ctyInts), nil
	case []float32:
		floats := val.([]float32)
		var ctyFloats []cty.Value
		for _, f := range floats {
			ctyFloats = append(ctyFloats, cty.NumberVal(big.NewFloat(float64(f))))
		}
		return cty.ListVal(ctyFloats), nil
	case []float64:
		floats := val.([]float64)
		var ctyFloats []cty.Value
		for _, f := range floats {
			ctyFloats = append(ctyFloats, cty.NumberVal(big.NewFloat(f)))
		}
		return cty.ListVal(ctyFloats), nil
	}
	return cty.NilVal, fmt.Errorf("unsupported type %s", reflect.TypeOf(val).String())
}
