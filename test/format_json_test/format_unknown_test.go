package format_json_test

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/reproio/terraform-j2md/internal/format"
	"reflect"
	"testing"
)

func TestFormatUnknownPlan(t *testing.T) {
	type args struct {
		old *tfjson.Plan
	}
	tests := []struct {
		name    string
		args    args
		want    *tfjson.Plan
		wantErr bool
	}{
		{
			name: "plain string",
			args: args{
				old: &tfjson.Plan{
					ResourceChanges: []*tfjson.ResourceChange{
						{
							Change: &tfjson.Change{
								Actions: tfjson.Actions{tfjson.ActionUpdate},
								Before: map[string]interface{}{
									"foo": "bar",
								},
								After: map[string]interface{}{},
								AfterUnknown: map[string]interface{}{
									"foo": true,
								},
							},
						},
					},
				},
			},
			want: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: tfjson.Actions{tfjson.ActionUpdate},
							Before: map[string]interface{}{
								"foo": "bar",
							},
							After: map[string]interface{}{
								"foo": "(known after apply)",
							},
							AfterUnknown: map[string]interface{}{
								"foo": true,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := format.FormatUnknownPlan(tt.args.old)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJsonPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatJsonPlan() got = \n%v\n, want \n%v", got.ResourceChanges[0].Change.After, tt.want.ResourceChanges[0].Change.After)
			}
		})
	}
}
