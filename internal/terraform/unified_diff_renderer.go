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
	return renderUnifiedDiff(r.ResourceChange, r.EnableEscapeHTML, 3, false)
}

// renderUnifiedDiff renders a resource change's before/after as a unified diff.
// context is the number of unchanged surrounding lines to show. Shared by the
// resource-change and drift renderers.
//
// trimTrailingNewline drops the newline json.Encode appends. Drift needs this:
// at high Context, SplitLines on a newline-terminated string yields a trailing
// empty element that surfaces as a spurious diff line. The change renderer must
// NOT trim — it would alter existing (upstream) change output for resources
// whose change lands within Context lines of the end.
func renderUnifiedDiff(rc *tfjson.ResourceChange, escapeHTML bool, context int, trimTrailingNewline bool) (string, error) {
	before, err := marshalDiffValue(rc.Change.Before, escapeHTML, trimTrailingNewline)
	if err != nil {
		return "", fmt.Errorf("invalid resource change (before): %w", err)
	}
	after, err := marshalDiffValue(rc.Change.After, escapeHTML, trimTrailingNewline)
	if err != nil {
		return "", fmt.Errorf("invalid resource change (after): %w", err)
	}
	// Try to parse JSON string in values
	replacer := strings.NewReplacer(`\n`, "\n  ", `\"`, "\"")
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(replacer.Replace(string(before))),
		B:       difflib.SplitLines(replacer.Replace(string(after))),
		Context: context,
	}
	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", fmt.Errorf("failed to create diff: %w", err)
	}

	return diffText, nil
}

func marshalDiffValue(v any, escapeHTML, trimTrailingNewline bool) ([]byte, error) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(escapeHTML)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	if trimTrailingNewline {
		return bytes.TrimRight(buffer.Bytes(), "\n"), nil
	}
	return buffer.Bytes(), nil
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
