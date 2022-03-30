### 2 to add, 0 to change, 2 to destroy.
- add
    - null_resource.aaa
    - null_resource.bbb
- destroy
    - null_resource.eee
    - null_resource.fff
<details><summary>Change details (Click me)</summary>

```diff
# null_resource.aaa will be created
@@ -1 +1,3 @@
-null
+{
+  "triggers": null
+}
```

```diff
# null_resource.bbb will be created
@@ -1 +1,3 @@
-null
+{
+  "triggers": null
+}
```

```diff
# null_resource.eee will be destroyed
@@ -1,4 +1 @@
-{
-  "id": "5480444040244548212",
-  "triggers": null
-}
+null
```

```diff
# null_resource.fff will be destroyed
@@ -1,4 +1 @@
-{
-  "id": "6136636772109947887",
-  "triggers": null
-}
+null
```

</details>
