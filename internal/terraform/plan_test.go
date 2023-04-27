package terraform

import (
	"bytes"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"os"
	"reflect"
	"testing"
)

func testDataPath(name, suffix string) string {
	return fmt.Sprintf("../../test/testdata/%s/%s", name, suffix)
}

func Test_newPlanData(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "no_changes", wantErr: false},
		{name: "single_add", wantErr: false},
		{name: "single_change", wantErr: false},
		{name: "single_destroy", wantErr: false},
		{name: "single_replace", wantErr: false},
		{name: "all_types_mixed", wantErr: false},
		{name: "aws_sample", wantErr: false},
		{name: "iam_policy", wantErr: false},
		{name: "invalid_json", wantErr: true},
		{name: "not_json", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFilePath := testDataPath(tt.name, "show.json")
			input, err := os.ReadFile(inputFilePath)
			if err != nil {
				t.Errorf("cannot open input file: %s", inputFilePath)
				return
			}

			_, err = NewPlanData(input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlanData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_render(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "no_changes", wantErr: false},
		{name: "single_add", wantErr: false},
		{name: "single_change", wantErr: false},
		{name: "single_destroy", wantErr: false},
		{name: "single_replace", wantErr: false},
		{name: "all_types_mixed", wantErr: false},
		{name: "aws_sample", wantErr: false},
		{name: "iam_policy", wantErr: false},
		{name: "include_code_fence", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFilePath := testDataPath(tt.name, "show.json")
			input, err := os.ReadFile(inputFilePath)
			if err != nil {
				t.Errorf("cannot open input file: %s", inputFilePath)
				return
			}

			plan, err := NewPlanData(input)
			if err != nil {
				t.Errorf("cannot parse JSON as plan: %v", err)
				return
			}

			got := bytes.Buffer{}
			err = plan.Render(&got)
			if (err != nil) != tt.wantErr {
				t.Errorf("render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			expectedFilePath := testDataPath(tt.name, "expected.md")
			expected, err := os.ReadFile(expectedFilePath)
			if err != nil {
				t.Errorf("cannot open expected file: %s", expectedFilePath)
				return
			}
			if got.String() != string(expected) {
				t.Errorf("render() = %v, want %v", got.String(), string(expected))
				return
			}
		})
	}
}

func TestResourceChangeData_hideSensitiveValue(t *testing.T) {
	type fields struct {
		ResourceChange *tfjson.ResourceChange
	}
	type args struct {
		change    interface{}
		sensitive interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "non-sensitive value",
			fields: fields{
				ResourceChange: &tfjson.ResourceChange{},
			},
			args: args{
				change: map[string]interface{}{
					"text": "dummy",
				},
				sensitive: false,
			},
			want: map[string]interface{}{
				"text": "dummy",
			},
			wantErr: false,
		},
		{
			name: "sensitive value",
			fields: fields{
				ResourceChange: &tfjson.ResourceChange{},
			},
			args: args{
				change: map[string]interface{}{
					"password": "dummy",
				},
				sensitive: map[string]interface{}{
					"password": true,
				},
			},
			want: map[string]interface{}{
				"password": "(sensitive value)",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceChangeData{
				ResourceChange: tt.fields.ResourceChange,
			}
			got, err := r.maskSensitiveValue(tt.args.change, tt.args.sensitive)
			if (err != nil) != tt.wantErr {
				t.Errorf("maskSensitiveValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("maskSensitiveValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
