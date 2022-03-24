### 2 to add, 0 to change, 2 to destroy.
- add
	- null_resource.aaa 
	- null_resource.bbb 
- destroy
	- null_resource.eee 
	- null_resource.fff 
<details><summary>Change details (Click me)</summary>

```diff
  any(
+ 	map[string]any{"triggers": nil},
  )

```
```diff
  any(
+ 	map[string]any{"triggers": nil},
  )

```
```diff
  any(
- 	map[string]any{"id": string("5480444040244548212"), "triggers": nil},
  )

```
```diff
  any(
- 	map[string]any{"id": string("6136636772109947887"), "triggers": nil},
  )

```
</details>
