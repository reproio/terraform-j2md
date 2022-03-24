package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
	//"github.com/pmezard/go-difflib/difflib"
	"github.com/google/go-cmp/cmp"
)
func getResourceNames(report []*tfjson.ResourceChange) {
	for _, i := range report {

		fmt.Printf("\t- %s \n", i.Address)
	}
	return
}
func getResourceDiff(report []*tfjson.ResourceChange) {
	for _, i := range report {
		if diff := cmp.Diff(i.Change.Before, i.Change.After); diff != "" {
			fmt.Printf("```diff\n%s\n```\n", diff)
		}
	}
	return
}
func main() {
	var plan tfjson.Plan

	err := json.NewDecoder(os.Stdin).Decode(&plan)
	if err != nil {
		log.Fatal(err)
		return
	}

	var report struct {
		Add     []*tfjson.ResourceChange
		Change  []*tfjson.ResourceChange
		Destroy []*tfjson.ResourceChange
	}

	//fmt.Println(plan.TerraformVersion)
	for _, c := range plan.ResourceChanges {
		if c.Change.Actions.NoOp() {
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
		//fmt.Println(c.Address)
	}

	// 2 to add, 0 to change, 2 to destroy.
	fmt.Printf("### %d to add, %d to change, %d to destroy.\n", len(report.Add), len(report.Change), len(report.Destroy))

	// リソース名を表示
	if report.Add != nil{
		fmt.Println("- add")
		getResourceNames(report.Add)
	}
	if report.Change != nil{
		fmt.Println("- change")
		getResourceNames(report.Change)
	}
	if report.Destroy != nil{
		fmt.Println("- destroy")
		getResourceNames(report.Destroy)
	};

	//展開して差分を表示する
	fmt.Println("<details><summary>Change details (Click me)</summary>\n")
	getResourceDiff(report.Add)
	getResourceDiff(report.Change)
	getResourceDiff(report.Destroy)
	fmt.Println("</details>")
}
