package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pmezard/go-difflib/difflib"
)
func getResourceNames(report []*tfjson.ResourceChange) string{
	var body string
	for _, i := range report {

		body += fmt.Sprintf("\t- %s\n", i.Address)
	}
	return body
}
func getResourceDiff(report *tfjson.ResourceChange) string{
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
	BeforeData, err := json.MarshalIndent(report.Change.Before, "", "  ")
	AfterData, err := json.MarshalIndent(report.Change.After, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	a := string(BeforeData)
	b := string(AfterData)
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
	}
	diffText, _ := difflib.GetUnifiedDiffString(diff)
	return fmt.Sprintf("\n```diff\n# %s.%s %s\n%s```\n", report.Type, report.Name, message, diffText)
}
func main() {
	var plan tfjson.Plan
	var body, diff string

	err := json.NewDecoder(os.Stdin).Decode(&plan)
	if err != nil {
		log.Fatal(err)
		return
	}

	var report struct {
		Add     []*tfjson.ResourceChange
		Change  []*tfjson.ResourceChange
		Destroy []*tfjson.ResourceChange
		Replace []*tfjson.ResourceChange
	}

	//fmt.Println(plan.TerraformVersion)
	for _, c := range plan.ResourceChanges {
		if c.Change.Actions.NoOp() {
			continue
		}
		if c.Change.Actions.Read() {
			continue
		}
		if c.Change.Actions.Create() {
			report.Add = append(report.Add, c)
		}
		if c.Change.Actions.Update() {
			report.Change = append(report.Change, c)
		}
		if c.Change.Actions.Delete() {
			report.Destroy = append(report.Destroy, c)
		}
		if c.Change.Actions.Replace() {
			report.Replace = append(report.Replace, c)
		}
		diff += getResourceDiff(c)
		
	}

	addCount := len(report.Add) + len(report.Replace)
	destroyCount := len(report.Destroy) + len(report.Replace)
	body += fmt.Sprintf("### %d to add, %d to change, %d to destroy.\n", addCount, len(report.Change), destroyCount)

	// リソース名を表示
	if report.Add != nil{
		body += fmt.Sprintln("- add")
		body += getResourceNames(report.Add)
	}
	if report.Change != nil{
		body += fmt.Sprintln("- change")
		body += getResourceNames(report.Change)
	}
	if report.Destroy != nil{
		body += fmt.Sprintln("- destroy")
		body += getResourceNames(report.Destroy)
	}
	if report.Replace != nil{
		body += fmt.Sprintln("- replace")
		body += getResourceNames(report.Replace)
	}

	//展開して差分を表示する
	body += fmt.Sprintf("<details><summary>Change details (Click me)</summary>\n%s\n</details>", diff)

	fmt.Printf("%s\n", body)
}
