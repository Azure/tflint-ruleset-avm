package attrvalue

import (
	"fmt"
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

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
