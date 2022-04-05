package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
)

const planTemplateBody = `### {{.CreatedCount}} to add, {{.UpdatedCount}} to change, {{.DeletedCount}} to destroy, {{.ReplacedCount}} to replace.
{{if .CreatedNames}}- add{{ range .CreatedNames }}
    - {{.Address -}}
{{end}}{{end -}}
{{if .UpdatedNames}}- change{{ range .UpdatedNames }}
    - {{.Address -}}
{{end}}{{end -}}
{{if .DeletedNames}}- delete{{ range .DeletedNames }}
    - {{.Address -}}
{{end}}{{end -}}
{{if .ReplacedNames}}- replace{{ range .ReplacedNames }}
    - {{.Address -}}
{{end}}{{end}}
{{if .ChangedResult}}<details><summary>Change details</summary>
{{ range .ChangedResult }}
{{.Backquote}}diff
# {{.ReportChanges.Type}}.{{.ReportChanges.Name}} {{.Message}}
{{.Diff}}{{.Backquote}}
{{end}}{{end}}
</details>
`

type PlanTemplate struct {
	CreatedCount  int
	UpdatedCount  int
	DeletedCount  int
	ReplacedCount int
	CreatedNames  []*tfjson.ResourceChange
	UpdatedNames  []*tfjson.ResourceChange
	DeletedNames  []*tfjson.ResourceChange
	ReplacedNames []*tfjson.ResourceChange
	ChangedResult []DiffTemplate
}
type DiffTemplate struct {
	ReportChanges *tfjson.ResourceChange
	Message       string
	Diff          string
	Backquote     string
}

func createResourceDiffString(resourceChanges *tfjson.ResourceChange) (string, error) {
	beforeData, err := json.MarshalIndent(resourceChanges.Change.Before, "", "  ")
	if err != nil {
		return "", fmt.Errorf("invalid resource changes (before): %w", err)
	}
	afterData, err := json.MarshalIndent(resourceChanges.Change.After, "", "  ")
	if err != nil {
		return "", fmt.Errorf("invalid resource changes (after) : %w", err)
	}
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(string(beforeData)),
		B:       difflib.SplitLines(string(afterData)),
		Context: 3,
	}
	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", fmt.Errorf("failed to create diff: %w", err)
	}
	return diffText, nil
}

func NewTemplateData(plan tfjson.Plan) (*PlanTemplate, error) {
	var report struct {
		Add     []*tfjson.ResourceChange
		Change  []*tfjson.ResourceChange
		Destroy []*tfjson.ResourceChange
		Replace []*tfjson.ResourceChange
	}
	diffs := []DiffTemplate{}

	for _, c := range plan.ResourceChanges {
		if c.Change.Actions.NoOp() || c.Change.Actions.Read() {
			continue
		}

		var message string
		switch {
		case c.Change.Actions.Create():
			report.Add = append(report.Add, c)
			message = "will be created"
		case c.Change.Actions.Update():
			report.Change = append(report.Change, c)
			message = "will be updated in-place"
		case c.Change.Actions.Delete():
			report.Destroy = append(report.Destroy, c)
			message = "will be destroyed"
		case c.Change.Actions.Replace():
			report.Replace = append(report.Replace, c)
			message = "will be replaced"
		}
		diff, err := createResourceDiffString(c)
		if err != nil {
			return nil, fmt.Errorf("invalid resource changes: %w", err)
		}
		diffs = append(diffs, DiffTemplate{
			ReportChanges: c,
			Message:       message,
			Diff:          diff,
			Backquote:     "```",
		})
	}

	return &PlanTemplate{
		CreatedCount:  len(report.Add),
		UpdatedCount:  len(report.Change),
		DeletedCount:  len(report.Destroy),
		ReplacedCount: len(report.Replace),
		CreatedNames:  report.Add,
		UpdatedNames:  report.Change,
		DeletedNames:  report.Destroy,
		ReplacedNames: report.Replace,
		ChangedResult: diffs,
	}, nil
}

func render(input string) (string, error) {
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

func run() int {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("cannot read stdin: %v", err)
		return 1
	}

	output, err := render(string(input))
	if err != nil {
		fmt.Printf("cannot convert: %v", err)
		return 1
	}

	fmt.Print(output)
	return 0
}

func main() {
	os.Exit(run())
}
