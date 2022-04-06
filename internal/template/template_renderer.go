package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
)

const planTemplateBody = `### {{.CreatedCount}} to add, {{.UpdatedCount}} to change, {{.DeletedCount}} to destroy, {{.ReplacedCount}} to replace.
{{- if .CreatedNames}}
- add{{ range .CreatedNames }}
    - {{.Address -}}
{{end}}{{end}}
{{- if .UpdatedNames}}
- change{{ range .UpdatedNames }}
    - {{.Address -}}
{{end}}{{end}}
{{- if .DeletedNames}}
- destroy{{ range .DeletedNames }}
    - {{.Address -}}
{{end}}{{end}}
{{- if .ReplacedNames}}
- replace{{ range .ReplacedNames }}
    - {{.Address -}}
{{end}}{{end}}
{{if .ChangedResult -}}
<details><summary>Change details</summary>
{{ range .ChangedResult }}
{{.Backquote}}diff
# {{.ReportChanges.Type}}.{{.ReportChanges.Name}} {{.Message}}
{{.Diff}}{{.Backquote}}
{{end}}
</details>
{{end}}`

func Render(input string) (string, error) {
	var plan tfjson.Plan
	err := json.Unmarshal([]byte(input), &plan)
	if err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	data, err := NewTemplateData(plan)
	if err != nil {
		return "", fmt.Errorf("invalid plan data: %w", err)
	}

	planTemplate, err := template.New("plan").Parse(planTemplateBody)
	if err != nil {
		return "", fmt.Errorf("invalid template text: %w", err)
	}

	var body bytes.Buffer
	if err := planTemplate.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return body.String(), nil
}
