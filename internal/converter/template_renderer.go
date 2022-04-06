package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
)

const planTemplateBody = `### {{.CreatedCount}} to add, {{.UpdatedCount}} to change, {{.DeletedCount}} to destroy, {{.ReplacedCount}} to replace.
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
{{.Backquote}}diff
# {{.ResourceChange.Type}}.{{.ResourceChange.Name}} {{.HeaderSuffix}}
{{.GetUnifiedDiffString}}{{.Backquote}}
{{end}}
</details>
{{end}}`

func Render(input string) (string, error) {
	var plan tfjson.Plan
	err := json.Unmarshal([]byte(input), &plan)
	if err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	planData, err := NewPlanData(plan)
	if err != nil {
		return "", fmt.Errorf("invalid plan data: %w", err)
	}

	planTemplate, err := template.New("plan").Parse(planTemplateBody)
	if err != nil {
		return "", fmt.Errorf("invalid template text: %w", err)
	}

	var output bytes.Buffer
	if err := planTemplate.Execute(&output, planData); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return output.String(), nil
}
