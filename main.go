package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
)

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

	fmt.Println(plan.TerraformVersion)
	for _, c := range plan.ResourceChanges {
		if c.Change.Actions.NoOp() {
			continue
		}

		if c.Change.Actions.Create() {
			report.Add = append(report.Add, c)
		}
		if c.Change.Actions.Delete() {
			report.Destroy = append(report.Destroy, c)
		}

		fmt.Println(c.Address)
	}

	fmt.Printf("%d to add, %d to change, %d to destroy.\n", len(report.Add), len(report.Change), len(report.Destroy))
}
