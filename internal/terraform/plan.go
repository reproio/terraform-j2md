package terraform

import (
	"encoding/json"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-json/sanitize"
	"github.com/reproio/terraform-j2md/internal/format"
	"io"
	"text/template"
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
{{- if .MovedAddresses}}
- moved{{ range .MovedAddresses }}
    - {{. -}}
{{end}}{{end}}
{{if .ResourceChanges -}}
<details><summary>Change details</summary>
{{ range .ResourceChanges }}
{{codeFence}}diff
# {{.Header}}
{{.Render}}{{codeFence}}
{{end}}
</details>
{{end}}`

type PlanData struct {
	Config            *Config
	CreatedAddresses  []string
	UpdatedAddresses  []string
	DeletedAddresses  []string
	ReplacedAddresses []string
	MovedAddresses    []string
	ResourceChanges   []ResourceChangeData
}

type Config struct {
	EscapeHTML bool
}

func NewPlanData(input io.Reader, escapeHTML bool) (*PlanData, error) {
	var err error
	var plan tfjson.Plan
	if err := json.NewDecoder(input).Decode(&plan); err != nil {
		return nil, fmt.Errorf("cannot parse input: %w", err)
	}

	processedPlan, err := processPlan(&plan)
	if err != nil {
		return nil, err
	}

	planData := PlanData{
		Config: &Config{
			EscapeHTML: escapeHTML,
		},
	}
	for _, c := range processedPlan.ResourceChanges {

		if isMovedBlock(c) {
			planData.MovedAddresses = append(planData.MovedAddresses, fmt.Sprintf("%s (from %s)", c.Address, c.PreviousAddress))
			planData.ResourceChanges = append(planData.ResourceChanges, ResourceChangeData{
				ResourceChange: c,
				Renderer:       NewMovedBlockRenderer(c),
			})
			continue
		}

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
			Renderer:       NewUnifiedDiffRenderer(c, escapeHTML),
		})
	}
	return &planData, nil
}

func processPlan(plan *tfjson.Plan) (*tfjson.Plan, error) {
	var err error

	for i := range plan.ResourceChanges {
		plan.ResourceChanges[i].Change, err = sanitize.SanitizeChange(plan.ResourceChanges[i].Change, sanitize.DefaultSensitiveValue)
		if err != nil {
			return nil, fmt.Errorf("failed to sanitize change: %w", err)
		}

		plan.ResourceChanges[i].Change, err = format.FormatJsonChange(plan.ResourceChanges[i].Change)
		if err != nil {
			return nil, fmt.Errorf("failed to format json change: %w", err)
		}

		plan.ResourceChanges[i].Change, err = format.FormatUnknownChange(plan.ResourceChanges[i].Change)
		if err != nil {
			return nil, fmt.Errorf("failed to format unknown change: %w", err)
		}
	}

	return plan, nil
}

func isMovedBlock(rc *tfjson.ResourceChange) bool {
	return rc.Change.Actions.NoOp() && rc.PreviousAddress != ""
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

type ResourceChangeData struct {
	ResourceChange *tfjson.ResourceChange
	Renderer       ResourceChangeDataRenderer
}

type ResourceChangeDataRenderer interface {
	Render() (string, error)
	Header() string
}

func (r ResourceChangeData) Render() (string, error) {
	return r.Renderer.Render()
}

func (r ResourceChangeData) Header() string {
	return r.Renderer.Header()
}
