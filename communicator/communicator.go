package communicator

import (
	"context"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
)

// A Communicator is the interface used to communicate with APIs
// that will eventually return metadata. Communicators
// allow you to get information from remote APi, databases, etc.
//
// Communicators must be safe for concurrency, meaning multiple calls to
// any method may be called at the same time.
type Communicator interface {
	// Configured
	Configured() bool

	// Configure Communicator
	Configure(map[string]interface{}) error
	// Fetch metadata from remote storage
	Fetch(context.Context, map[string]interface{}) ([]map[string]interface{}, error)
}

// GetCommunicator returns Communicator by type.
func GetCommunicator(communicator_type string) (Communicator, error) {
	switch communicator_type {
	case "http", "HTTP":
		return NewRestCommunicator(), nil
	default:
		return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableCommunicator, communicator_type)
	}
}

// GetSectionCommunicator returns communicator from configuration file.
// By default http communicator will be used.
// Example YAML config for `section` that will return new `RestCommunicator`:
//     section:
//         communicator:
//             type: "HTTP"
func GetSectionCommunicator(section string) (Communicator, error) {
	communicator_type := config.GetStringDefault(fmt.Sprintf("%s.%s.type", section, config.CFG_PREFIX_COMMUNICATOR), "http")
	switch communicator_type {
	case "http", "HTTP":
		comm := NewRestCommunicator()
		var cfg_params map[string]interface{}
		cfg_params = config.ConvertMapStringToInterface(
			config.GetStringMapStringTemplated(section, config.CFG_PREFIX_COMMUNICATOR))
		if _, ok := cfg_params["section"]; !ok {
			cfg_params["section"] = fmt.Sprintf("%s.%s.type", section, config.CFG_PREFIX_COMMUNICATOR)
		}
		if _, ok := cfg_params["param"]; !ok {
			cfg_params["param"] = config.CFG_COMMUNICATOR_PARAMS_KEY
		}

		if err := comm.Configure(cfg_params); err != nil {
			return nil, err
		}

		return comm, nil
	default:
		return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableCommunicator, communicator_type)
	}
}
