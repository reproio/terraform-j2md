package terraform

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
)

const planTemplateBody = `### {{len .CreatedAddresses}} to add, {{len .UpdatedAddresses}} to change, {{len .DeletedAddresses}} to destroy, {{len .ReplacedAddresses}} to replace.
{{- if .CreatedAddresses}}
- add{{ range .CreatedAddresses }}
    - {{. -}}
{{end}}{{end}}
{{- if .UpdatedAddresses}}
- change{{ range .UpdatedAddresses }}
    - {{. -}}
{{end}}{{end}}
{{- if .DeletedAddresses}}
- destroy{{ range .DeletedAddresses }}
    - {{. -}}
{{end}}{{end}}
{{- if .ReplacedAddresses}}
- replace{{ range .ReplacedAddresses }}
    - {{. -}}
{{end}}{{end}}
{{if .ResourceChanges -}}
<details><summary>Change details</summary>
{{ range .ResourceChanges }}
{{codeFence}}diff
# {{.ResourceChange.Type}}.{{.ResourceChange.Name}} {{.HeaderSuffix}}
{{.GetUnifiedDiffString}}{{codeFence}}
{{end}}
</details>
{{end}}`

type PlanData struct {
	CreatedAddresses  []string
	UpdatedAddresses  []string
	DeletedAddresses  []string
	ReplacedAddresses []string
	ResourceChanges   []ResourceChangeData
}
type ResourceChangeData struct {
	ResourceChange *tfjson.ResourceChange
}

func (r ResourceChangeData) renderChange(change interface{}, sensitive interface{}) ([]byte, error) {
	re, err := r.maskSensitiveValue(change, sensitive)
	if err != nil {
		return nil, fmt.Errorf("maskSensitiveValue fails: %w", err)
	}
	return json.MarshalIndent(re, "", "  ")
}

func hasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

func (r ResourceChangeData) maskSensitiveValue(change interface{}, sensitive interface{}) (interface{}, error) {
	switch sensitive.(type) {
	case bool:
		if !sensitive.(bool) {
			return change, nil
		}
	case map[string]interface{}:
		mappedSensitive := sensitive.(map[string]interface{})
		if len(mappedSensitive) == 0 {
			return change, nil
		}

		var re = map[string]interface{}{}
		for k, v := range change.(map[string]interface{}) {
			if hasKey(mappedSensitive, k) && reflect.ValueOf(mappedSensitive[k]).Type().String() == "bool" && mappedSensitive[k].(bool) && v != "" {
				re[k] = "(sensitive value)"
			} else {
				re[k] = v
			}
		}
		return re, nil
	default:
		return nil, fmt.Errorf("renderChange can't handling type of sensitive: %T", sensitive)
	}
	return nil, nil
}

func (r ResourceChangeData) GetUnifiedDiffString() (string, error) {
	before, err := r.renderChange(r.ResourceChange.Change.Before, r.ResourceChange.Change.BeforeSensitive)
	if err != nil {
		return "", fmt.Errorf("invalid resource changes (before): %w", err)
	}
	after, err := r.renderChange(r.ResourceChange.Change.After, r.ResourceChange.Change.AfterSensitive)
	if err != nil {
		return "", fmt.Errorf("invalid resource changes (after) : %w", err)
	}
	// Try to parse JSON string in values
	replacer := strings.NewReplacer(`\n`, "\n  ", `\"`, "\"")
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(replacer.Replace(string(before))),
		B:       difflib.SplitLines(replacer.Replace(string(after))),
		Context: 3,
	}
	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", fmt.Errorf("failed to create diff: %w", err)
	}

	return diffText, nil
}

func (r ResourceChangeData) HeaderSuffix() string {
	switch {
	case r.ResourceChange.Change.Actions.Create():
		return "will be created"
	case r.ResourceChange.Change.Actions.Update():
		return "will be updated in-place"
	case r.ResourceChange.Change.Actions.Delete():
		return "will be destroyed"
	case r.ResourceChange.Change.Actions.Replace():
		return "will be replaced"
	}
	return ""
}

func (plan *PlanData) Render(w io.Writer) error {
	funcMap := template.FuncMap{
		"codeFence": func() string {
			return "````````"
		},
	}
	planTemplate, err := template.New("plan").Funcs(funcMap).Parse(planTemplateBody)
	if err != nil {
		return fmt.Errorf("invalid template text: %w", err)
	}

	if err := planTemplate.Execute(w, plan); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}

func NewPlanData(input []byte) (*PlanData, error) {
	var plan tfjson.Plan
	if err := json.Unmarshal(input, &plan); err != nil {
		return nil, fmt.Errorf("cannot parse input: %w", err)
	}
	planData := PlanData{}
	for _, c := range plan.ResourceChanges {
		if c.Change.Actions.NoOp() || c.Change.Actions.Read() {
			continue
		}

		switch {
		case c.Change.Actions.Create():
			planData.CreatedAddresses = append(planData.CreatedAddresses, c.Address)
		case c.Change.Actions.Update():
			planData.UpdatedAddresses = append(planData.UpdatedAddresses, c.Address)
		case c.Change.Actions.Delete():
			planData.DeletedAddresses = append(planData.DeletedAddresses, c.Address)
		case c.Change.Actions.Replace():
			planData.ReplacedAddresses = append(planData.ReplacedAddresses, c.Address)
		}
		planData.ResourceChanges = append(planData.ResourceChanges, ResourceChangeData{
			ResourceChange: c,
		})
	}
	return &planData, nil
}
