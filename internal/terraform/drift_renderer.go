package terraform

import (
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

// driftDiffContext is the number of unchanged surrounding lines shown for drift
// diffs. Drift output lives inside a collapsed <details> block, so we show the
// full resource state (large value) rather than the tight window used for
// planned changes.
const driftDiffContext = 1000

type DriftRenderer struct {
	ResourceChange   *tfjson.ResourceChange
	EnableEscapeHTML bool
}

func NewDriftRenderer(resourceChange *tfjson.ResourceChange, enableEscapeHTML bool) *DriftRenderer {
	return &DriftRenderer{ResourceChange: resourceChange, EnableEscapeHTML: enableEscapeHTML}
}

func (r *DriftRenderer) Render() (string, error) {
	return renderUnifiedDiff(r.ResourceChange, r.EnableEscapeHTML, driftDiffContext, true)
}

func (r *DriftRenderer) Header() string {
	return fmt.Sprintf("%s has drifted from state", r.ResourceChange.Address)
}
