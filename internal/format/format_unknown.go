package format

import (
	"errors"
	tfjson "github.com/hashicorp/terraform-json"
)

func FormatUnknownPlan(plan *tfjson.Plan) (*tfjson.Plan, error) {
	var err error
	if plan == nil {
		return nil, errors.New("nil plan supplied")
	}

	for i := range plan.ResourceChanges {
		plan.ResourceChanges[i].Change, err = formatUnknownChange(plan.ResourceChanges[i].Change)
		if err != nil {
			return nil, err
		}
	}

	return plan, nil
}

func formatUnknownChange(change *tfjson.Change) (*tfjson.Change, error) {
	if change.Actions.Update() {
		for k, v := range change.AfterUnknown.(map[string]interface{}) {
			switch v.(type) {
			case bool:
				change.After.(map[string]interface{})[k] = "(known after apply)"
			}
		}
	}
	return change, nil
}
