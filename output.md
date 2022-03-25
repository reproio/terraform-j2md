### 2 to add, 0 to change, 2 to destroy.
- add
	- null_resource.aaa 
	- null_resource.bbb 
- destroy
	- null_resource.eee 
	- null_resource.fff 
<details><summary>Change details (Click me)</summary>

```diff
resource null_resource aaa
- null
+ {
+   "triggers": null
+ }
```
```diff
resource null_resource bbb
- null
+ {
+   "triggers": null
+ }
```
```diff
resource null_resource eee
- {
-   "id": "5480444040244548212",
-   "triggers": null
- }
+ null
```
```diff
resource null_resource fff
- {
-   "id": "6136636772109947887",
-   "triggers": null
- }
+ null
```
</details>
