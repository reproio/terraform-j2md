package render

import (
	"bytes"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"text/template"
)

const movedBlockTemplateBody = `resource "{{.ResourceChange.Type}}" "{{.ResourceChange.Name}}" {
{{.Attributes -}}
}
`

// These attributes are important (https://github.com/hashicorp/terraform/blob/v1.5.6/internal/command/jsonformat/computed/renderers/block.go#L19-L23)
var importantAttributes = []string{
	"id",
	"name",
	"tags",
}

type MovedBlockRenderer struct {
	ResourceChange *tfjson.ResourceChange
}

func NewMovedBlockRenderer(resourceChange *tfjson.ResourceChange) *MovedBlockRenderer {
	return &MovedBlockRenderer{ResourceChange: resourceChange}
}

func (r *MovedBlockRenderer) Render() (string, error) {
	var buff bytes.Buffer
	t, err := template.New("plan").Parse(movedBlockTemplateBody)
	if err != nil {
		return "", fmt.Errorf("invalid template text: %w", err)
	}

	if err := t.Execute(&buff, r); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return buff.String(), nil
}

func (r *MovedBlockRenderer) Header() string {
	return fmt.Sprintf("%s has moved to %s", r.ResourceChange.PreviousAddress, r.ResourceChange.Address)
}

func (r *MovedBlockRenderer) Attributes() string {
	var buff bytes.Buffer
	for _, attr := range importantAttributes {
		if v, ok := r.ResourceChange.Change.After.(map[string]interface{})[attr]; ok {
			buff.WriteString(fmt.Sprintf("  %-*s = %s\n", 2, attr, r.value(v)))
		}
	}
	return buff.String()
}

func (r *MovedBlockRenderer) value(v any) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
