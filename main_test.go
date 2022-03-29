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
	}

	defer input.Close()
	os.Stdin = input
	main()
}

func ExampleNoChanges() {
	RunTestMain("test_data/no_changes/show.json")
	// Output:
	// 0 to add, 0 to change, 0 to destroy.
}

func ExampleSingleAdd() {
	RunTestMain("test_data/single_add/show.json")
	// Output:
	// ### 1 to add, 0 to change, 0 to destroy.
	// - add
	// 	- null_resource.foo
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// resource null_resource foo
	// @@ -1 +1,3 @@
	// -null
	// +{
	// +  "triggers": null
	// +}
	// ```
	//
	// </details>
}

func ExampleSingleDestroy() {
	RunTestMain("test_data/single_destroy/show.json")
	// Output:
	// ### 0 to add, 0 to change, 1 to destroy.
	// - destroy
	// 	- null_resource.foo
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// resource null_resource foo
	// @@ -1,4 +1 @@
	// -{
	// -  "id": "317876227733854172",
	// -  "triggers": null
	// -}
	// +null
	// ```
	//
	// </details>
}

func ExampleSingleChange() {
	RunTestMain("test_data/single_change/show.json")
	// Output:
	// 	### 0 to add, 1 to change, 0 to destroy.
	// - change
	//	- env_variable.test1
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// resource env_variable test1
	// @@ -1,5 +1,5 @@
	//  {
	//    "id": "test1",
	// -  "name": "test1",
	// +  "name": "test1_changed",
	//    "value": ""
	//  }
	// ```
	//
	// </details>
}

func ExampleReplaceChange() {
	RunTestMain("test_data/single_replace/show.json")
	// Output:
	// ### 1 to add, 0 to change, 1 to destroy.
	// - replace
	// 	- random_id.test
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// resource random_id test
	// @@ -1,10 +1,5 @@
	//  {
	// -  "b64_std": "qddo6VPNl1g=",
	// -  "b64_url": "qddo6VPNl1g",
	// -  "byte_length": 8,
	// -  "dec": "12238365863745263448",
	// -  "hex": "a9d768e953cd9758",
	// -  "id": "qddo6VPNl1g",
	// +  "byte_length": 10,
	//    "keepers": null,
	//    "prefix": null
	//  }
	// ```
	//
	// </details>
}

func ExampleAllMixed() {
	RunTestMain("test_data/all_mixed/show.json")
	// Output:
	// ### 2 to add, 1 to change, 2 to destroy.
	// - add
	// 	- env_variable.test5
	// - change
	// 	- env_variable.test2
	// - destroy
	// 	- env_variable.test3
	// - replace
	// 	- random_id.test4
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// resource env_variable test2
	// @@ -1,5 +1,5 @@
	//  {
	//    "id": "test2",
	// -  "name": "test2",
	// +  "name": "test2_changed",
	//    "value": ""
	//  }
	// ```
	//
	// ```diff
	// resource env_variable test3
	// @@ -1,5 +1 @@
	// -{
	// -  "id": "test3",
	// -  "name": "test3",
	// -  "value": ""
	// -}
	// +null
	// ```
	//
	// ```diff
	// resource env_variable test5
	// @@ -1 +1,3 @@
	// -null
	// +{
	// +  "name": "test5"
	// +}
	// ```
	//
	// ```diff
	// resource random_id test4
	// @@ -1,10 +1,5 @@
	//  {
	// -  "b64_std": "m6S5W82/OFA=",
	// -  "b64_url": "m6S5W82_OFA",
	// -  "byte_length": 8,
	// -  "dec": "11215292776004401232",
	// -  "hex": "9ba4b95bcdbf3850",
	// -  "id": "m6S5W82_OFA",
	// +  "byte_length": 10,
	//    "keepers": null,
	//    "prefix": null
	//  }
	// ```
	//
	// </details>
}
