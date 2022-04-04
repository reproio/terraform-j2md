package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestAll(t *testing.T) {
	type args struct {
		dataPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "No Changes", args: args{dataPath: "testdata/no_changes"}},
		{name: "Single Create", args: args{dataPath: "testdata/single_add"}},
		{name: "Single Change", args: args{dataPath: "testdata/single_change"}},
		{name: "Single Destroy", args: args{dataPath: "testdata/single_destroy"}},
		{name: "Single Replace", args: args{dataPath: "testdata/single_replace"}},
		{name: "All Change Types Mixed", args: args{dataPath: "testdata/all_mixed"}},
		{name: "AWS Resource Changes", args: args{dataPath: "testdata/aws_mixed"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataPath := tt.args.dataPath
			expectedFilePath := dataPath + "/expected.md"
			expected, err := os.ReadFile(expectedFilePath)
			if err != nil {
				t.Errorf("cannot open input file: %s", expectedFilePath)
			}

			actual := RunCase(t, dataPath+"/show.json")
			if string(expected) != actual {
				diff := difflib.UnifiedDiff{
					A:       difflib.SplitLines(actual),
					B:       difflib.SplitLines(string(expected)),
					Context: 3,
				}
				diffText, err := difflib.GetUnifiedDiffString(diff)
				if err != nil {
					t.Error("result does not matched, and could not render diff")
				} else {
					t.Errorf("result does not matched:\n %s", diffText)
				}
			}
		})
	}
}

func RunCase(t *testing.T, inputFile string) string {
	input, err := os.Open(inputFile)
	if err != nil {
		t.Errorf("cannot open input file: %s", inputFile)
	}
	defer input.Close()

	r, w, err := os.Pipe()
	if err != nil {
		t.Error("cannot exec os.Pipe")
	}
	defer func() {
		w.Close()
		r.Close()
	}()

	originStdin := os.Stdin
	originStdout := os.Stdout
	os.Stdin = input
	os.Stdout = w
	defer func() {
		os.Stdin = originStdin
		os.Stdout = originStdout
	}()
	main()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
