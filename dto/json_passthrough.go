package dto

import "encoding/json"

func collectExtraFields(data []byte, knownKeys ...string) map[string]json.RawMessage {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil
	}
	for _, key := range knownKeys {
		delete(fields, key)
	}
	if len(fields) == 0 {
		return nil
	}
	return fields
}

func marshalWithExtraFields(value any, extraFields map[string]json.RawMessage) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(extraFields) == 0 {
		return data, nil
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	for key, value := range extraFields {
		if _, exists := fields[key]; exists {
			continue
		}
		fields[key] = value
	}
	return json.Marshal(fields)
}
