package utils

import (
	cast "github.com/spf13/cast"
	"strconv"
	"time"
)

func GetFirstStringFromMap(v map[string]interface{}, keys []string) (string, bool) {
	for _, k := range keys {
		if _, ok := v[k]; ok {

			if val, ok1 := v[k].(string); ok1 {

				return val, true
			}

		}
	}
	return "", false
}

func GetFirstTimeFromMap(v map[string]interface{}, keys []string) (time.Time, bool) {
	for _, k := range keys {
		if _, ok := v[k]; ok {
			if val, err := cast.ToTimeE(v[k]); err == nil {
				return val, true
			}

			switch t := v[k].(type) {
			case time.Time:
				return t, true
			case string:
				if i, err := strconv.Atoi(t); err == nil {
					return time.Unix(int64(i), 0), true
				}
			case int:
				return time.Unix(int64(t), 0), true
			case int64:
				return time.Unix(t, 0), true
			case float64:
				return time.Unix(int64(int(t)), 0), true
			}
		}
	}
	var tempTime time.Time
	return tempTime, false
}
