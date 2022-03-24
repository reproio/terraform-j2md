package main

import (
	"log"
	"os"
)

func RunTestMain(inputFile string) {
	originStdin := os.Stdin
	defer func() {
		os.Stdin = originStdin
	}()

	input, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("cannot open input file: %s", inputFile)
		return
	}

	defer input.Close()
	os.Stdin = input
	main()
}

func ExampleNoChanges() {
	RunTestMain("test_data/empty.json")
	// Output:
	// 0 to add, 0 to change, 0 to destroy.
}

func ExampleAdd() {
	RunTestMain("test_data/add.json")
	// Output:
	// 1 to add, 0 to change, 0 to destroy.
	//
	// - add:
	//   - null_resource.foo
	//
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// + resource "null_resource" "foo" {
	// +   id = (known after apply)
	// +  }
	// ```
	//
	// </details>
}

func ExampleDestroy() {
	RunTestMain("test_data/destroy.json")
	// Output:
	// 0 to add, 0 to change, 1 to destroy.
	//
	// - destroy:
	//   - null_resource.foo
	//
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// - resource "null_resource" "foo" {
	// -   id = "317876227733854172"
	// -  }
	// ```
	//
	// </details>
}
