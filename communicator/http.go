package communicator

import (
	"bytes"
	"fmt"
	"strings"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	config "github.com/weldpua2008/suprasched/config"
)

// RestCommunicator represents API communicator.
type RestCommunicator struct {
	Communicator
	section    string
	param      string
	method     string
	url        string
	headers    map[string]string
	configured bool
}

// NewRestCommunicator prepare struct communicator for HTTP requests
func NewRestCommunicator() *RestCommunicator {
	return &RestCommunicator{}
}

// Configured checks if RestCommunicator is configured.
func (s *RestCommunicator) Configured() bool {
	return s.configured
}

// Configure reads configuration propertoes from global configuration and
// from argument.
func (s *RestCommunicator) Configure(params map[string]interface{}) error {

	if _, ok := params["section"]; ok {
		s.section = params["section"].(string)
	}
	if _, ok := params["param"]; ok {
		s.param = params["param"].(string)
	}
	s.method = "POST"
	if _, ok := params["method"]; ok {
		s.method = strings.ToUpper(params["method"].(string))
	}
	if _, ok := params["url"]; ok {
		s.url = params["url"].(string)
	}
	s.headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	if _, ok := params["headers"]; ok {
		if v, ok1 := params["headers"].(map[string]string); ok1 {
			for k1, v1 := range v {
				s.headers[k1] = v1
			}
		}

	}
	s.configured = true
	// log.Tracef("%v", s.method)
	return nil
}

func (s *RestCommunicator) Fetch(ctx context.Context, params map[string]interface{}) (result []map[string]interface{}, err error) {
	var jsonStr []byte

	all_params := config.GetStringMapStringTemplated(s.section, s.param)
	for k, v := range params {
		if v1, ok := v.(string); ok {
			all_params[k] = v1
		}
	}
	var req *http.Request
	var rawResponse map[string]interface{}

	if len(all_params) > 0 {
		jsonStr, err = json.Marshal(&all_params)
		if err != nil {
			log.Tracef("\nFailed to marshal request %s  to %s \nwith %s\n", s.method, s.url, jsonStr)
			return nil, fmt.Errorf("Failed to marshal request due %s", err)
		}

		req, err = http.NewRequest(s.method, s.url, bytes.NewBuffer(jsonStr))
	} else {
		req, err = http.NewRequest(s.method, s.url, nil)
	}
	if err != nil {
		return result, fmt.Errorf("Failed to create request due %s", err)
	}
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Duration(15 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request due %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read response body got %s", err)
	}
	if (resp.StatusCode > 202) || (resp.StatusCode < 200) {
		log.Tracef("\nMaking request %s  to %s \nwith %s\nStatusCode %d Response %s\n", s.method, s.url, jsonStr, resp.StatusCode, body)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, fmt.Errorf("error Unmarshal response: %s due %s", body, err)
		}
		result = append(result, rawResponse)
	}

	return result, nil

}
