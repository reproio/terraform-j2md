package format_json_test

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-json/sanitize"
	"github.com/reproio/terraform-j2md/internal/format"
	"reflect"
	"testing"
)

func TestSanitize(t *testing.T) {
	type args struct {
		old *tfjson.Change
	}
	tests := []struct {
		name    string
		args    args
		want    *tfjson.Change
		wantErr bool
	}{
		{
			name: "001",
			args: args{
				old: &tfjson.Change{
					Before: map[string]interface{}{
						"test1": "bar",
					},
					After: map[string]interface{}{
						"test1": "bar",
					},
					BeforeSensitive: map[string]interface{}{
						"test1": true,
					},
					AfterSensitive: map[string]interface{}{
						"test1": true,
					},
				},
			},
			want: &tfjson.Change{
				Before: map[string]interface{}{
					"test1": sanitize.DefaultSensitiveValue,
				},
				After: map[string]interface{}{
					"test1": sanitize.DefaultSensitiveValue,
				},
				BeforeSensitive: map[string]interface{}{
					"test1": true,
				},
				AfterSensitive: map[string]interface{}{
					"test1": true,
				},
			},
			wantErr: false,
		},
		{
			name: "002",
			args: args{
				old: &tfjson.Change{
					Before: map[string]interface{}{
						"test1": "bar",
					},
					After: map[string]interface{}{
						"test1": "bar",
					},
					//BeforeSensitive: map[string]interface{}{
					//	"test1": true,
					//},
					AfterSensitive: map[string]interface{}{
						"test1": true,
					},
				},
			},
			want: &tfjson.Change{
				Before: map[string]interface{}{
					"test1": "bar",
				},
				After: map[string]interface{}{
					"test1": sanitize.DefaultSensitiveValue,
				},
				AfterSensitive: map[string]interface{}{
					"test1": true,
				},
				//BeforeSensitive: map[string]interface{}{
				//	"test1": true,
				//},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitize.SanitizeChange(tt.args.old, sanitize.DefaultSensitiveValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sanitize() got = \n%v\n, want \n%v", got, tt.want)
			}
		})
	}
}

func TestFormatJsonPlan(t *testing.T) {
	type args struct {
		old *tfjson.Change
	}
	tests := []struct {
		name    string
		args    args
		want    *tfjson.Change
		wantErr bool
	}{
		{
			name: "plain string",
			args: args{
				old: &tfjson.Change{
					Before: "plain string",
					After:  "plain string",
				},
			},
			want: &tfjson.Change{
				Before: "plain string",
				After:  "plain string",
			},
			wantErr: false,
		},
		{
			name: "json string",
			args: args{
				old: &tfjson.Change{
					Before: `{"foo":"bar"}`,
					After:  `{"foo":"bar"}`,
				},
			},
			want: &tfjson.Change{
				Before: `{
  "foo": "bar"
}`,
				After: `{
  "foo": "bar"
}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := format.FormatJsonChange(tt.args.old)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJsonPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatJsonPlan() got = \n%v\n, want \n%v", got.After, tt.want.After)
			}
		})
	}
}
