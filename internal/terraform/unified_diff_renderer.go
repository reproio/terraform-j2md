package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
	"strings"
)

type UnifiedDiffRenderer struct {
	ResourceChange   *tfjson.ResourceChange
	EnableEscapeHTML bool
}

func NewUnifiedDiffRenderer(resourceChange *tfjson.ResourceChange, enableEscapeHTML bool) *UnifiedDiffRenderer {
	return &UnifiedDiffRenderer{ResourceChange: resourceChange, EnableEscapeHTML: enableEscapeHTML}
}

func (r *UnifiedDiffRenderer) Render() (string, error) {
	before, err := r.marshalChangeBefore()
	if err != nil {
		return "", fmt.Errorf("invalid resource changes (before): %w", err)
	}
	after, err := r.marshalChangeAfter()
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

func (r *UnifiedDiffRenderer) Header() string {
	header := fmt.Sprintf("%s %s", r.ResourceChange.Address, r.headerSuffix())

	return header
}

func (r *UnifiedDiffRenderer) headerSuffix() string {
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

func (r *UnifiedDiffRenderer) marshalChangeBefore() ([]byte, error) {
	return r.marshalChange(r.ResourceChange.Change.Before)
}

func (r *UnifiedDiffRenderer) marshalChangeAfter() ([]byte, error) {
	return r.marshalChange(r.ResourceChange.Change.After)
}

func (r *UnifiedDiffRenderer) marshalChange(v any) ([]byte, error) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(r.EnableEscapeHTML)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
