### 1 to add, 1 to change, 1 to destroy, 0 to replace.
- add
    - random_id.test2
- change
    - env_variable.test
- destroy
    - random_id.test
<details><summary>Change details</summary>

````````diff
# env_variable.test will be updated in-place
@@ -1,6 +1,6 @@
 {
   "id": "07ec30c9d869e0f6392f",
-  "name": "07ec30c9d869e0f6392f",
+  "name": "(known after apply)",
   "value": "REDACTED_SENSITIVE"
 }
 
````````

````````diff
# random_id.test will be destroyed
@@ -1,11 +1,2 @@
-{
-  "b64_std": "B+wwydhp4PY5Lw==",
-  "b64_url": "B-wwydhp4PY5Lw",
-  "byte_length": 10,
-  "dec": "37413512560416367458607",
-  "hex": "07ec30c9d869e0f6392f",
-  "id": "B-wwydhp4PY5Lw",
-  "keepers": null,
-  "prefix": null
-}
+null
 
````````

````````diff
# random_id.test2 will be created
@@ -1,2 +1,6 @@
-null
+{
+  "byte_length": 10,
+  "keepers": null,
+  "prefix": null
+}
 
````````

</details>
