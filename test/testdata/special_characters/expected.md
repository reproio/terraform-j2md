### 1 to add, 0 to change, 0 to destroy, 0 to replace.
- add
    - aws_instance.web
<details><summary>Change details</summary>

````````diff
# aws_instance.web will be created
@@ -1,2 +1,23 @@
-null
+{
+  "ami": "ami-04fc53a873660e525",
+  "credit_specification": [],
+  "get_password_data": false,
+  "hibernation": null,
+  "instance_type": "t3.micro",
+  "launch_template": [],
+  "source_dest_check": true,
+  "tags": {
+    "tag1": ">",
+    "tag2": "<",
+    "tag3": "&"
+  },
+  "tags_all": {
+    "tag1": ">",
+    "tag2": "<",
+    "tag3": "&"
+  },
+  "timeouts": null,
+  "user_data_replace_on_change": false,
+  "volume_tags": null
+}
 
````````

</details>
