### 2 to add, 0 to change, 0 to destroy, 0 to replace.
- add
    - aws_instance.web["t3.micro"]
    - aws_instance.web["t3.small"]
<details><summary>Change details</summary>

````````diff
# aws_instance.web["t3.micro"] will be created
@@ -1,2 +1,14 @@
-null
+{
+  "ami": "ami-04fc53a873660e525",
+  "credit_specification": [],
+  "get_password_data": false,
+  "hibernation": null,
+  "instance_type": "t3.micro",
+  "launch_template": [],
+  "source_dest_check": true,
+  "tags": null,
+  "timeouts": null,
+  "user_data_replace_on_change": false,
+  "volume_tags": null
+}
 
````````

````````diff
# aws_instance.web["t3.small"] will be created
@@ -1,2 +1,14 @@
-null
+{
+  "ami": "ami-04fc53a873660e525",
+  "credit_specification": [],
+  "get_password_data": false,
+  "hibernation": null,
+  "instance_type": "t3.small",
+  "launch_template": [],
+  "source_dest_check": true,
+  "tags": null,
+  "timeouts": null,
+  "user_data_replace_on_change": false,
+  "volume_tags": null
+}
 
````````

</details>
