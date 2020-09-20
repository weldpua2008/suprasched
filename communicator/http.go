package communicator

import (
	"bytes"
	"fmt"
	"strings"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	config "github.com/weldpua2008/suprasched/config"
)

func init() {
	Constructors[ConstructorsTypeRest] = TypeSpec{
		instance:    NewRestCommunicator,
		constructor: NewConfiguredRestCommunicator,
		Summary: `
RestCommunicator is the default implementation of Communicator and is
used by Default.
Underline it is used for HTTP connections only.`,
		Description: `
It supports the following params:
- ` + "`method`" + ` http method more [in this document](https://golang.org/pkg/net/http/#pkg-constants)
- ` + "`headers`" + ` To make a request with custom headers
- ` + "`url`" + ` to send in a request for the given URL.`,
	}
}

// RestCommunicator represents HTTP communicator.
type RestCommunicator struct {
	Communicator
	section    string
	param      string
	method     string
	url        string
	headers    map[string]string
	configured bool
	mu         sync.RWMutex
}

// NewRestCommunicator prepare struct communicator for HTTP requests
func NewRestCommunicator() Communicator {
	return &RestCommunicator{}
}

// NewRestCommunicator prepare struct communicator for HTTP requests
func NewConfiguredRestCommunicator(section string) (Communicator, error) {
	comm := NewRestCommunicator()
	var cfg_params map[string]interface{}
	cfg_params = config.ConvertMapStringToInterface(
		config.GetStringMapStringTemplated(section, config.CFG_PREFIX_COMMUNICATOR))
	if _, ok := cfg_params["section"]; !ok {
		cfg_params["section"] = fmt.Sprintf("%s.%s", section, config.CFG_PREFIX_COMMUNICATOR)
	}
	if _, ok := cfg_params["param"]; !ok {
		cfg_params["param"] = config.CFG_COMMUNICATOR_PARAMS_KEY
	}

	if err := comm.Configure(cfg_params); err != nil {
		return nil, err
	}

	return comm, nil
}

// Configured checks if RestCommunicator is configured.
func (s *RestCommunicator) Configured() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.configured
}

// Configure reads configuration propertoes from global configuration and
// from argument.
func (s *RestCommunicator) Configure(params map[string]interface{}) error {
    // log.Warningf("Configure %v", params)
	s.mu.Lock()
	defer s.mu.Unlock()

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
	return nil
}

func (s *RestCommunicator) Fetch(ctx context.Context, params map[string]interface{}) (result []map[string]interface{}, err error) {
	var jsonStr []byte
	var req *http.Request
	var rawResponse map[string]interface{}


	// all_params := config.GetStringMapStringTemplated(s.section, s.param)
    from:= map[string]string{
        "ClientId": config.C.ClientId,
        "ClusterId": config.C.ClusterId,
        "ClusterPool": config.C.ClusterPool,
        "ConfigVersion": config.C.ConfigVersion,
    }
    for k, v := range params {
        if v1, ok := v.(string); ok {
            from[k] = v1
        }
    }
    s.mu.RLock()
	defer s.mu.RUnlock()

    all_params := config.GetStringMapStringTemplatedFromMap(s.section, s.param, from)
    // log.Infof("\nall_params %v\ns.section %v , %v, \nfrom: %v", all_params, s.section, s.param,from)

	for k, v := range params {
		if v1, ok := v.(string); ok {
			all_params[k] = v1
		}
	}

	if len(all_params) > 0 {
		jsonStr, err = json.Marshal(&all_params)
		if err != nil {
			log.Tracef("\nFailed to marshal request %s  to %s \nwith %s\n", s.method, s.url, jsonStr)
			return nil, fmt.Errorf("%w due %s", ErrFailedMarshalRequest, err)
		}

		req, err = http.NewRequest(s.method, s.url, bytes.NewBuffer(jsonStr))
	} else {
		req, err = http.NewRequest(s.method, s.url, nil)
	}
	if err != nil {
		return result, fmt.Errorf("%w due %s", ErrFailedSendRequest, err)
	}
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Duration(15 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w got %s", ErrFailedSendRequest, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w got %s", ErrFailedReadResponseBody, err)
	}
	if (resp.StatusCode > 202) || (resp.StatusCode < 200) {

        log.Tracef("\nMaking request %s  to %s \nwith %s\nStatusCode %d Response %s\n", s.method, s.url, jsonStr, resp.StatusCode, body)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, fmt.Errorf("%w from body %s got %s", ErrFailedUnmarshalResponse, body, err)
		}
		result = append(result, rawResponse)
	}

	return result, nil

}
