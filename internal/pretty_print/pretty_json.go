package pretty_print

import (
	"encoding/json"
	"errors"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/jinzhu/copier"
)

func PrettyPrintPlan(old *tfjson.Plan) (*tfjson.Plan, error) {
	if old == nil {
		return nil, errors.New("nil plan supplied")
	}

	result, err := copyPlan(old)
	if err != nil {
		return nil, err
	}

	for i := range result.ResourceChanges {
		result.ResourceChanges[i].Change, err = prettyChange(result.ResourceChanges[i].Change)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func prettyChange(old *tfjson.Change) (*tfjson.Change, error) {
	result, err := copyChange(old)
	if err != nil {
		return nil, err
	}

	result.Before, err = prettyChangeValue(result.Before)
	if err != nil {
		return nil, err
	}
	result.After, err = prettyChangeValue(result.After)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func prettyChangeValue(old interface{}) (interface{}, error) {
	switch x := old.(type) {
	case []interface{}:
		for i, v := range x {
			result, err := prettyChangeValue(v)
			if err != nil {
				return nil, err
			}
			x[i] = result
		}
	case map[string]interface{}:
		for k, v := range x {
			result, err := prettyChangeValue(v)
			if err != nil {
				return nil, err
			}
			x[k] = result
		}
	case string:
		var j json.RawMessage
		if err := json.Unmarshal([]byte(old.(string)), &j); err == nil && json.Valid([]byte(old.(string))) {
			a, err := json.MarshalIndent(j, "", "  ")
			if err != nil {
				return "", err
			}
			return string(a), nil
		}
	}

	return old, nil
}

func copyPlan(old *tfjson.Plan) (*tfjson.Plan, error) {
	result := &tfjson.Plan{}
	err := copier.CopyWithOption(result, old, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func copyChange(old *tfjson.Change) (*tfjson.Change, error) {
	result := &tfjson.Change{}
	err := copier.CopyWithOption(result, old, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	return result, nil
}
