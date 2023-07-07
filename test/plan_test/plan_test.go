package plan_test

import (
	"bytes"
	"fmt"
	"github.com/reproio/terraform-j2md/internal/terraform"
	"os"
	"testing"
)

func testDataPath(name, suffix string) string {
	return fmt.Sprintf("../testdata/%s/%s", name, suffix)
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
			file, err := os.Open(inputFilePath)
			if err != nil {
				t.Errorf("cannot open input file: %s", inputFilePath)
				return
			}
			defer file.Close()

			_, err = terraform.NewPlanData(file)
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
		{name: "include_module", wantErr: false},
		{name: "special_characters", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFilePath := testDataPath(tt.name, "show.json")
			file, err := os.Open(inputFilePath)
			if err != nil {
				t.Errorf("cannot open input file: %s", inputFilePath)
				return
			}
			defer file.Close()

			plan, err := terraform.NewPlanData(file)
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
