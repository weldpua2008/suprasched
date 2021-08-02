package communicator

import (
	"context"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
	"strings"
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
func GetCommunicator(communicatorType string) (Communicator, error) {
	k := strings.ToUpper(communicatorType)
	if typeStruct, ok := Constructors[k]; ok {
		if comm := typeStruct.instance(); comm != nil {
			return comm, nil
		} else {
			return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableCommunicator, communicatorType)
		}
	}

	return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableCommunicator, communicatorType)
}

// GetSectionCommunicator returns communicator from configuration file.
// By default http communicator will be used.
// Example YAML config for `section` that will return new `RestCommunicator`:
//     section:
//         communicator:
//             type: "HTTP"
func GetSectionCommunicator(section string) (Communicator, error) {
	communicatorType := config.GetStringDefault(fmt.Sprintf("%s.%s.type", section, config.CFG_PREFIX_COMMUNICATOR), "http")
	k := strings.ToUpper(communicatorType)
	if typeStruct, ok := Constructors[k]; ok {
		if comm, err := typeStruct.constructor(section); err == nil {
			return comm, nil
		} else {
			return nil, err
		}

	}
	return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableCommunicator, communicatorType)
}

// GetCommunicatorsFromSection returns multiple communicators from configuration file.
// By default http communicator will be used.
// Example YAML config for `section` that will return new `RestCommunicator`:
//     section:
//         communicators:
//             my_communicator:
//                 type: "HTTP"
//             -:
//                 type: "HTTP"
func GetCommunicatorsFromSection(section string) ([]Communicator, error) {
	def := make(map[string]string)

	// comms := config.GetSliceStringMapStringTemplatedDefault(section, config.CFG_PREFIX_COMMUNICATORS, def)
	comms := config.GetMapStringMapStringTemplatedDefault(section, config.CFG_PREFIX_COMMUNICATORS, def)

	res := make([]Communicator, 0)
	for section, comm := range comms {
		if comm == nil {
			continue
		}
		communicatorType := ConstructorsTypeRest
		if commType, ok := comm["type"]; ok {
			communicatorType = commType
		}
		comm["section"] = section
		if _, ok := comm["param"]; !ok {
			comm["param"] = config.CFG_COMMUNICATOR_PARAMS_KEY
		}

		k := strings.ToUpper(communicatorType)
		if typeStruct, ok := Constructors[k]; ok {
			// if comm, err := type_struct.constructor(section); err == nil {
			communicatorInstance := typeStruct.instance()
			// log.Warningf("comm %v", comm)
			if err1 := communicatorInstance.Configure(config.ConvertMapStringToInterface(comm)); err1 != nil {
				log.Tracef("Can't configure %v communicator, got %v", communicatorType, comm)
				return nil, err1
			}
			res = append(res, communicatorInstance)
		}
	}
	if len(res) > 0 {
		return res, nil
	}
	return nil, fmt.Errorf("%w in section %s.\n", ErrNoSuitableCommunicator, section)
}
