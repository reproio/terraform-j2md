package converter

import (
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
)

type PlanData struct {
	CreatedCount      int
	UpdatedCount      int
	DeletedCount      int
	ReplacedCount     int
	CreatedAddresses  []string
	UpdatedAddresses  []string
	DeletedAddresses  []string
	ReplacedAddresses []string
	ResourceChanges   []ResourceChangeData
}
type ResourceChangeData struct {
	ResourceChange *tfjson.ResourceChange
	HeaderSuffix   string
}

func (ResourceChangeData) Backquote() string {
	return "```"
}

func (d ResourceChangeData) GetUnifiedDiffString() (string, error) {
	beforeData, err := json.MarshalIndent(d.ResourceChange.Change.Before, "", "  ")
	if err != nil {
		return "", fmt.Errorf("invalid resource changes (before): %w", err)
	}
	afterData, err := json.MarshalIndent(d.ResourceChange.Change.After, "", "  ")
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

func NewPlanData(plan tfjson.Plan) (*PlanData, error) {
	var report struct {
		Add     []string
		Change  []string
		Destroy []string
		Replace []string
	}
	diffs := []ResourceChangeData{}

	for _, c := range plan.ResourceChanges {
		if c.Change.Actions.NoOp() || c.Change.Actions.Read() {
			continue
		}

		var headerSuffix string
		switch {
		case c.Change.Actions.Create():
			report.Add = append(report.Add, c.Address)
			headerSuffix = "will be created"
		case c.Change.Actions.Update():
			report.Change = append(report.Change, c.Address)
			headerSuffix = "will be updated in-place"
		case c.Change.Actions.Delete():
			report.Destroy = append(report.Destroy, c.Address)
			headerSuffix = "will be destroyed"
		case c.Change.Actions.Replace():
			report.Replace = append(report.Replace, c.Address)
			headerSuffix = "will be replaced"
		}
		diffs = append(diffs, ResourceChangeData{
			ResourceChange: c,
			HeaderSuffix:   headerSuffix,
		})
	}

	return &PlanData{
		CreatedCount:      len(report.Add),
		UpdatedCount:      len(report.Change),
		DeletedCount:      len(report.Destroy),
		ReplacedCount:     len(report.Replace),
		CreatedAddresses:  report.Add,
		UpdatedAddresses:  report.Change,
		DeletedAddresses:  report.Destroy,
		ReplacedAddresses: report.Replace,
		ResourceChanges:   diffs,
	}, nil
}
