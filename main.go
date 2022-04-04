package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
)

const tmplStr = `### {{.CreatedCount}} to add, {{.UpdatedCount}} to change, {{.DeletedCount}} to destroy, {{.ReplacedCount}} to replace.
{{if .CreatedNames}}- add{{ range .CreatedNames }}
    - {{.Address -}}
{{ end }}{{end}}
{{if .UpdatedNames}}- change{{ range .UpdatedNames }}
    - {{.Address -}}
{{ end }}{{end}}
{{if .DeletedNames}}- delete{{ range .DeletedNames }}
    - {{.Address -}}
{{ end }}{{end}}
{{if .ReplacedNames}}- replace{{ range .ReplacedNames }}
    - {{.Address -}}
{{ end }}{{end}}
{{if .ChangedResult}}<details><summary>Change details</summary>
{{ range .ChangedResult }}
{{.Backquote}}diff
# {{.ReportChanges.Type}}.{{.ReportChanges.Name}} {{.Message}}
{{.Diff}}
{{.Backquote}}
{{end}}{{end}}
</details>
`

type CommonTemplate struct {
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

func createResourceDiffString(resourceChanges *tfjson.ResourceChange) (string, string, error) {
	var message string
	switch {
	case resourceChanges.Change.Actions.Create():
		message = "will be created"
	case resourceChanges.Change.Actions.Update():
		message = "will be updated in-place"
	case resourceChanges.Change.Actions.Delete():
		message = "will be destroyed"
	case resourceChanges.Change.Actions.Replace():
		message = "will be replaced"
	}
	beforeData, err := json.MarshalIndent(resourceChanges.Change.Before, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("error message : %w", err)
	}
	afterData, err := json.MarshalIndent(resourceChanges.Change.After, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("error message : %w", err)
	}
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(string(beforeData)),
		B:       difflib.SplitLines(string(afterData)),
		Context: 3,
	}
	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", "", fmt.Errorf("error message : %w", err)
	}
	return message, diffText, nil
}
func render(input string) (string, error) {
	var plan tfjson.Plan
	err := json.Unmarshal([]byte(input), &plan)
	if err != nil {
		return "", fmt.Errorf("error message : %w", err)
	}

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

		switch {
		case c.Change.Actions.Create():
			report.Add = append(report.Add, c)
		case c.Change.Actions.Update():
			report.Change = append(report.Change, c)
		case c.Change.Actions.Delete():
			report.Destroy = append(report.Destroy, c)
		case c.Change.Actions.Replace():
			report.Replace = append(report.Replace, c)
		}
		message, diff, err := createResourceDiffString(c)
		if err != nil {
			return "", fmt.Errorf("error message : %w", err)
		}
		diffData := DiffTemplate{
			ReportChanges: c,
			Message:       message,
			Diff:          diff,
			Backquote:     "```",
		}
		diffs = append(diffs, diffData)
	}

	data := CommonTemplate{
		CreatedCount:  len(report.Add),
		UpdatedCount:  len(report.Change),
		DeletedCount:  len(report.Destroy),
		ReplacedCount: len(report.Replace),
		CreatedNames:  report.Add,
		UpdatedNames:  report.Change,
		DeletedNames:  report.Destroy,
		ReplacedNames: report.Replace,
		ChangedResult: diffs,
	}

	tmpl, err := template.New("test_template").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("error message : %w", err)
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		return "", fmt.Errorf("error message : %w", err)
	}
	return "", nil
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
