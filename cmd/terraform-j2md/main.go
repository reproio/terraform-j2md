package main

import (
	"fmt"
	"io"
	"os"
	"terraform-j2md/internal/template"
)

func main() {
	os.Exit(run())
}

func run() int {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("cannot read stdin: %v", err)
		return 1
	}

	output, err := template.Render(string(input))
	if err != nil {
		fmt.Printf("cannot convert: %v", err)
		return 1
	}

	fmt.Print(output)
	return 0
}
