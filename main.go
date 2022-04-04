package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
)

func listResourceNames(header string, report []*tfjson.ResourceChange) string {
	if report == nil {
		return ""
	}

	body := "- " + header + "\n"
	for _, c := range report {
		body += fmt.Sprintf("    - %s\n", c.Address)
	}
	return body
}
func createResourceDiffString(report *tfjson.ResourceChange) string {
	var message string
	switch {
	case report.Change.Actions.Create():
		message = "will be created"
	case report.Change.Actions.Update():
		message = "will be updated in-place"
	case report.Change.Actions.Delete():
		message = "will be destroyed"
	case report.Change.Actions.Replace():
		message = "will be replaced"
	}
	beforeData, err := json.MarshalIndent(report.Change.Before, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	afterData, err := json.MarshalIndent(report.Change.After, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	a := string(beforeData)
	b := string(afterData)
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
	}
	diffText, _ := difflib.GetUnifiedDiffString(diff)
	return fmt.Sprintf("\n```diff\n# %s.%s %s\n%s```\n", report.Type, report.Name, message, diffText)
}

func render(input string) (string, error) {
	var plan tfjson.Plan
	err := json.Unmarshal([]byte(input), &plan)
	if err != nil {
		return "", fmt.Errorf("input format is invalid: %w", err)
	}

	var report struct {
		Add     []*tfjson.ResourceChange
		Change  []*tfjson.ResourceChange
		Destroy []*tfjson.ResourceChange
		Replace []*tfjson.ResourceChange
	}

	var diff string
	for _, c := range plan.ResourceChanges {
		if c.Change.Actions.NoOp() || c.Change.Actions.Read() {
			continue
		}

		switch {
		case c.Change.Actions.Create():
			report.Add = append(report.Add, c)
		case c.Change.Actions.Update():
			report.Change = append(report.Change, c)
		case c.Change.Actions.Delete():
			report.Destroy = append(report.Destroy, c)
		case c.Change.Actions.Replace():
			report.Replace = append(report.Replace, c)
		}
		diff += createResourceDiffString(c)
	}

	var body string
	body += fmt.Sprintf("### %d to add, %d to change, %d to destroy.\n", len(report.Add)+len(report.Replace), len(report.Change), len(report.Destroy)+len(report.Replace))
	body += listResourceNames("add", report.Add)
	body += listResourceNames("change", report.Change)
	body += listResourceNames("destroy", report.Destroy)
	body += listResourceNames("replace", report.Replace)
	if len(diff) != 0 {
		body += fmt.Sprintf("<details><summary>Change details (Click me)</summary>\n%s\n</details>\n", diff)
	}

	return body, nil
}

func run() int {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("cannot read stdin: %v", err)
		return 1
	}

	output, err := render(string(input))
	if err != nil {
		fmt.Printf("cannot convert: %v", err)
		return 1
	}

	fmt.Print(output)
	return 0
}

func main() {
	os.Exit(run())
}
