package format

import (
	tfjson "github.com/hashicorp/terraform-json"
)

func FormatUnknownChange(change *tfjson.Change) (*tfjson.Change, error) {
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
