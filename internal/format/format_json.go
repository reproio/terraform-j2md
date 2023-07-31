package format

import (
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
)

func FormatJsonChange(change *tfjson.Change) (*tfjson.Change, error) {
	var err error

	change.Before, err = formatJsonChangeValue(change.Before)
	if err != nil {
		return nil, err
	}
	change.After, err = formatJsonChangeValue(change.After)
	if err != nil {
		return nil, err
	}

	return change, nil
}

func formatJsonChangeValue(old interface{}) (interface{}, error) {
	switch x := old.(type) {
	case []interface{}:
		for i, v := range x {
			result, err := formatJsonChangeValue(v)
			if err != nil {
				return nil, err
			}
			x[i] = result
		}
	case map[string]interface{}:
		for k, v := range x {
			result, err := formatJsonChangeValue(v)
			if err != nil {
				return nil, err
			}
			x[k] = result
		}
	case string:
		var j json.RawMessage
		if json.Valid([]byte(old.(string))) && json.Unmarshal([]byte(old.(string)), &j) == nil {
			a, err := json.MarshalIndent(j, "", "  ")
			if err != nil {
				return "", err
			}
			return string(a), nil
		}
	}

	return old, nil
}
