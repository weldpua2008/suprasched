package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGetFirstStringFromMap(t *testing.T) {
	cases := []struct {
		in     map[string]interface{}
		keys   []string
		result string
		want   bool
	}{
		{
			in: map[string]interface{}{
				"search": "found",
				"wrong":  1,
				"extra":  make(chan int),
			},
			keys:   []string{"Search", "search"},
			result: "found",
			want:   true,
		},
		{
			in: map[string]interface{}{
				"search": "found",
				"wrong":  1,
				"extra":  make(chan int),
			},
			keys:   []string{"wrong", "extra"},
			result: "",
			want:   false,
		},
	}

	for _, tc := range cases {
		result, got := GetFirstStringFromMap(tc.in, tc.keys)
		if tc.want != got {
			t.Errorf("want %v, got %v", tc.want, got)
		}
		if tc.result != result {
			t.Errorf("want %v, got %v", tc.result, result)
		}
	}
}

func TestGetFirstTimeFromMap(t *testing.T) {
	var not_found time.Time
	now := time.Now()
	cases := []struct {
		in     map[string]interface{}
		keys   []string
		result time.Time
		want   bool
	}{
		{
			in: map[string]interface{}{
				"search": now,
				"wrong":  1,
				"extra":  make(chan int),
			},
			keys:   []string{"Search", "search"},
			result: now,
			want:   true,
		},
		{
			in: map[string]interface{}{
				"search": "found",
				"wrong":  make(chan int),
				"extra":  make(chan int),
			},
			keys:   []string{"wrong", "extra"},
			result: not_found,
			want:   false,
		},
		{
			in: map[string]interface{}{
				"search": "found",
				"wrong":  make(chan int),
				"extra":  make(chan int),
			},
			keys:   []string{"wrong", "extra"},
			result: not_found,
			want:   false,
		},
		{
			in: map[string]interface{}{
				"search": fmt.Sprintf("%v", now.Unix()),
				"wrong":  1,
				"extra":  make(chan int),
			},
			keys:   []string{"Search", "search"},
			result: time.Unix(now.Unix(), 0),
			want:   true,
		},
		{
			in: map[string]interface{}{
				"search": int(now.Unix()),
				"wrong":  1,
				"extra":  make(chan int),
			},
			keys:   []string{"Search", "search"},
			result: time.Unix(now.Unix(), 0),
			want:   true,
		},
		{
			in: map[string]interface{}{
				"search": int64(now.Unix()),
				"wrong":  1,
				"extra":  make(chan int),
			},
			keys:   []string{"Search", "search"},
			result: time.Unix(now.Unix(), 0),
			want:   true,
		},
		{
			in: map[string]interface{}{
				"search": float64(now.Unix()),
				"wrong":  1,
				"extra":  make(chan int),
			},
			keys:   []string{"Search", "search"},
			result: time.Unix(now.Unix(), 0),
			want:   true,
		},
	}

	for n, tc := range cases {
		result, got := GetFirstTimeFromMap(tc.in, tc.keys)
		if tc.want != got {
			t.Errorf("want %v, got %v", tc.want, got)
		}
		if tc.result != result {
			t.Errorf("want %v result %v, got result %v, tc.in %v", n, tc.result, result, tc.in)
		}
	}
}
