package main

import (
	"fmt"
	"os"

	"github.com/reproio/terraform-j2md/internal/terraform"
)

func main() {
	os.Exit(run())
}

func run() int {
	planData, err := terraform.NewPlanData(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot parse input as Terraform plan JSON: %v", err)
		return 1
	}
	if err = planData.Render(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "cannot render: %v", err)
		return 1
	}
	return 0
}
