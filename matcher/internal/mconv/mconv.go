package mconv

import (
	"fmt"
	"strconv"
)

func ConvToFloat64(v any) (float64, error) {
	switch e := v.(type) {
	case string:
		return strconv.ParseFloat(e, 64)
	case float64:
		return e, nil
	case float32:
		return float64(e), nil
	case int:
		return float64(e), nil
	case int32:
		return float64(e), nil
	case int64:
		return float64(e), nil
	default:
		return 0, fmt.Errorf("value cannot be parsed to float64")
	}
}
