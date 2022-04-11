package terraform

import (
	"bytes"
	"os"
	"testing"
)

func Test_newPlanData(t *testing.T) {
	type args struct {
		testDataPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "No Changes", args: args{testDataPath: "../../test/testdata/no_changes"}, wantErr: false},
		{name: "Single Create", args: args{testDataPath: "../../test/testdata/single_add"}, wantErr: false},
		{name: "Single Change", args: args{testDataPath: "../../test/testdata/single_change"}, wantErr: false},
		{name: "Single Destroy", args: args{testDataPath: "../../test/testdata/single_destroy"}, wantErr: false},
		{name: "Single Replace", args: args{testDataPath: "../../test/testdata/single_replace"}, wantErr: false},
		{name: "All Change Types Mixed", args: args{testDataPath: "../../test/testdata/all_mixed"}, wantErr: false},
		{name: "AWS Resource Changes", args: args{testDataPath: "../../test/testdata/aws_mixed"}, wantErr: false},
		{name: "Invalid JSON input", args: args{testDataPath: "../../test/testdata/invalid_json"}, wantErr: true},
		{name: "Invalid format input", args: args{testDataPath: "../../test/testdata/not_json"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFilePath := tt.args.testDataPath + "/show.json"
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
	type args struct {
		testDataPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "No Changes", args: args{testDataPath: "../../test/testdata/no_changes"}, wantErr: false},
		{name: "Single Create", args: args{testDataPath: "../../test/testdata/single_add"}, wantErr: false},
		{name: "Single Change", args: args{testDataPath: "../../test/testdata/single_change"}, wantErr: false},
		{name: "Single Destroy", args: args{testDataPath: "../../test/testdata/single_destroy"}, wantErr: false},
		{name: "Single Replace", args: args{testDataPath: "../../test/testdata/single_replace"}, wantErr: false},
		{name: "All Change Types Mixed", args: args{testDataPath: "../../test/testdata/all_mixed"}, wantErr: false},
		{name: "AWS Resource Changes", args: args{testDataPath: "../../test/testdata/aws_mixed"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFilePath := tt.args.testDataPath + "/show.json"
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

			expectedFilePath := tt.args.testDataPath + "/expected.md"
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
