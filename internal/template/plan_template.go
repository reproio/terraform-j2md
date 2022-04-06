package template

import (
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
)

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
}

func (DiffTemplate) Backquote() string {
	return "```"
}

func (d DiffTemplate) Diff() (string, error) {
	beforeData, err := json.MarshalIndent(d.ReportChanges.Change.Before, "", "  ")
	if err != nil {
		return "", fmt.Errorf("invalid resource changes (before): %w", err)
	}
	afterData, err := json.MarshalIndent(d.ReportChanges.Change.After, "", "  ")
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
		diffs = append(diffs, DiffTemplate{
			ReportChanges: c,
			Message:       message,
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
