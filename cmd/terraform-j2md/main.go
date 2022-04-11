package main

import (
	"fmt"
	"io"
	"os"
	"terraform-j2md/internal/terraform"
)

func main() {
	os.Exit(run())
}

func run() int {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read stdin: %v", err)
		return 1
	}
	planData, err := terraform.NewPlanData(input)
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
