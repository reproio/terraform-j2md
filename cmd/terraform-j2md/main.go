package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/reproio/terraform-j2md/internal/terraform"
)

var (
	escapeHTML = true
)

func main() {
	noEscapeHTML := flag.Bool("no-escape-html", false, "prevent <, >, and & from being escaped in JSON strings")
	flag.Parse()
	if *noEscapeHTML {
		escapeHTML = false
	}
	os.Exit(run())
}

func run() int {
	planData, err := terraform.NewPlanData(os.Stdin, escapeHTML)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot parse input as Terraform plan JSON: %v", err)
		return 1
	}
	if err = planData.Render(os.Stdout, escapeHTML); err != nil {
		fmt.Fprintf(os.Stderr, "cannot render: %v", err)
		return 1
	}
	return 0
}
