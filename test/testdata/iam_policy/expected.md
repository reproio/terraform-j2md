### 0 to add, 1 to change, 0 to destroy, 0 to replace.
- change
    - aws_iam_policy.test_policy
<details><summary>Change details</summary>

````````diff
# aws_iam_policy.test_policy will be updated in-place
@@ -12,7 +12,8 @@
       "Action": [
         "autoscaling:Describe*",
         "ec2:Describe*",
-        "elasticloadbalancing:Describe*"
+        "elasticloadbalancing:Describe*",
+        "health:Describe*"
       ],
       "Resource": "*"
     }
````````

</details>
