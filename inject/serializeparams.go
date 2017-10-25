package inject

import (
	"encoding/json"

	api "github.com/antha-lang/antha/api/v1"
	any "github.com/golang/protobuf/ptypes/any"
)

// SerializeValueToElementParameters is a helper function to serialize
// inject values into api messages.
func SerializeValueToElementParameters(v Value) api.ElementParameters {
	m := make(map[string]*any.Any, len(v))

	for key, thing := range v {
		bs, err := json.Marshal(thing)

		if err != nil {
			panic(err.Error())
		}

		m[key] = &any.Any{Value: bs}
	}

	return api.ElementParameters{Parameters: m}
}
