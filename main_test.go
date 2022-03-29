package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	cases := map[string]string{
		"No Changes":             "test_data/no_changes",
		"Single Create":          "test_data/single_add",
		"Single Change":          "test_data/single_change",
		"Single Destroy":         "test_data/single_destroy",
		"Single Replace":         "test_data/single_replace",
		"All Change Types Mixed": "test_data/all_mixed",
		"AWS Resource Changes":   "test_data/aws_mixed",
	}

	for caseName, casePath := range cases {
		t.Run(caseName, func(t *testing.T) {
			expectedFilePath := casePath + "/expected.md"
			expected, err := os.ReadFile(expectedFilePath)
			if err != nil {
				t.Errorf("cannot open input file: %s", expectedFilePath)
			}

			actual := RunCase(t, casePath+"/show.json")
			if string(expected) != actual {
				t.Errorf("\ngot:\n%v\nwant:\n%v", actual, string(expected))
			}
		})
	}
}

func RunCase(t *testing.T, inputFile string) string {
	originStdin := os.Stdin
	originStdout := os.Stdout
	defer func() {
		os.Stdin = originStdin
		os.Stdout = originStdout
	}()

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

	os.Stdout = w
	os.Stdin = input
	main()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
