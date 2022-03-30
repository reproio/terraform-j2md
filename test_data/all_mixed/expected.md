### 2 to add, 1 to change, 2 to destroy.
- add
    - env_variable.test5
- change
    - env_variable.test2
- destroy
    - env_variable.test3
- replace
    - random_id.test4
<details><summary>Change details (Click me)</summary>

```diff
# env_variable.test2 will be updated in-place
@@ -1,5 +1,5 @@
 {
   "id": "test2",
-  "name": "test2",
+  "name": "test2_changed",
   "value": ""
 }
```

```diff
# env_variable.test3 will be destroyed
@@ -1,5 +1 @@
-{
-  "id": "test3",
-  "name": "test3",
-  "value": ""
-}
+null
```

```diff
# env_variable.test5 will be created
@@ -1 +1,3 @@
-null
+{
+  "name": "test5"
+}
```

```diff
# random_id.test4 will be replaced
@@ -1,10 +1,5 @@
 {
-  "b64_std": "m6S5W82/OFA=",
-  "b64_url": "m6S5W82_OFA",
-  "byte_length": 8,
-  "dec": "11215292776004401232",
-  "hex": "9ba4b95bcdbf3850",
-  "id": "m6S5W82_OFA",
+  "byte_length": 10,
   "keepers": null,
   "prefix": null
 }
```

</details>
