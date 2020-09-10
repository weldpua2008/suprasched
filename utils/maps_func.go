package utils

import (
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

			switch t := v[k].(type) {
			case string:
				if i, err := strconv.Atoi(t); err == nil {
					return time.Unix(int64(i), 0), true

				}
			case int:
				return time.Unix(int64(t), 0), true

			case float64:
				return time.Unix(int64(int(t)), 0), true
			}
		}
	}
	var temp_time time.Time
	return temp_time, false
}
