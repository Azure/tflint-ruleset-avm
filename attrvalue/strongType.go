package attrvalue

func toStrongTypePtr(val any) any {
	switch t := val.(type) {
	case *string:
		return t
	case *int:
		return t
	case *float64:
		return t
	case *bool:
		return t
	case string:
		return &t
	case int:
		return &t
	case float64:
		return &t
	case bool:
		return &t
	case *[]int32:
		return t
	case []int32:
		return &t
	}
	return nil
}
