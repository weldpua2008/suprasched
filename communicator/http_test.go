package communicator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	Number int    `json:"number"`
	Str    string `json:"str"`
}

func TestFetch(t *testing.T) {
	config.LoadCfgForTests(t, "fixtures/fetch_http.yml")
	var globalGot string
	responses := []Response{
		{
			Number: 1,
			Str:    "Str",
		},
		{
			Number: 2,
			Str:    "Str1",
		},
	}
	// Response server.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var c Response
		if len(responses) > 1 {
			c, responses = responses[0], responses[1:]
		} else if len(responses) == 1 {
			c = responses[0]
		}
		c1 := make([]Response, 0)
		c1 = append(c1, c)
		js, err := json.Marshal(&c1)
		if err != nil {
			log.Tracef("Failed to marshal for '%v' due %v", c, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, errWrite := w.Write(js); errWrite != nil {
			t.Errorf("Can't w.Write %v due %v\n", js, err)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ReadAll %s", err)
		}
		globalGot = string(fmt.Sprintf("%s", b))
	}))
	defer func() {
		srv.Close()
	}()

	cases := []struct {
		section      string
		params       map[string]interface{}
		want         map[string]interface{}
		wantErr      error
		wantFetchErr error
	}{
		{
			section: "get",
			params: map[string]interface{}{
				"method": "GET",
				"url":    srv.URL,
			},
			want: map[string]interface{}{
				"Number": responses[0].Number,
				"Str":    responses[0].Str,
			},
			wantErr:      nil,
			wantFetchErr: nil,
		},
		{
			section: "get",
			params: map[string]interface{}{
				"method": "GET",
				"url":    srv.URL,
			},
			want: map[string]interface{}{
				"Number": responses[1].Number,
				"Str":    responses[1].Str,
			},
			wantErr:      nil,
			wantFetchErr: nil,
		},
	}
	for _, tc := range cases {
		result, got := GetSectionCommunicator(tc.section)
		if (tc.wantErr == nil) && (tc.wantErr != got) {
			t.Errorf("want %v, got %v", tc.wantErr, got)
		} else if (tc.want == nil) && (!result.Configured()) {
			t.Errorf("want %v, got %v, res %v", true, result.Configured(), result)
		} else {
			if !errors.Is(got, tc.wantErr) {
				t.Errorf("want %v, got %v, res %v", tc.wantErr, got, result)
			}
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // cancel when we are getting the kill signal or exit
		sent_params := map[string]interface{}{
			"a": "a",
		}
		if err := result.Configure(tc.params); err != nil {
			t.Errorf("want %v, got %v, results %v", nil, err, result)
		}
		c, _ := result.(*RestCommunicator)

		if len(c.url) < 1 {
			t.Errorf("want url len, got %v", result)

		}

		ret, getFetchErr := result.Fetch(ctx, sent_params)
		if (tc.wantFetchErr == nil) && (tc.wantFetchErr != getFetchErr) {
			t.Errorf("want %v, got %v", tc.wantFetchErr, getFetchErr)
		// WARNING: Keys are always in Lower case.
		} else if (tc.wantFetchErr == nil) && (ret[0]["str"] != tc.want["Str"]) {
			t.Errorf("want %v, got %v", tc.want["Str"], ret[0]["str"])
		} else {
			if !errors.Is(getFetchErr, tc.wantFetchErr) {
				t.Errorf("want %v, got %v, res %v", tc.wantFetchErr, getFetchErr, result)
			}
		}
		if len(globalGot) < 1 {
			t.Errorf("want len > 0 , got %v, send params %v", globalGot, sent_params)
		}

	}
}
