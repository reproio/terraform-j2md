### 0 to add, 0 to change, 0 to destroy, 0 to replace.

### ⚠️ Drift Detected (2 resources)
- aws_instance.web
- aws_security_group.allow_http
<details><summary>Drift details</summary>

````````diff
# aws_instance.web has drifted from state
@@ -1,9 +1,10 @@
 {
   "ami": "ami-12345678",
   "arn": "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
   "id": "i-1234567890abcdef0",
   "instance_type": "t2.micro",
   "tags": {
+    "Environment": "production",
     "Name": "web-server"
   }
 }
````````

````````diff
# aws_security_group.allow_http has drifted from state
@@ -1,15 +1,23 @@
 {
   "description": "Allow HTTP traffic",
   "id": "sg-0123456789abcdef0",
   "ingress": [
     {
       "cidr_blocks": [
         "0.0.0.0/0"
       ],
       "from_port": 80,
       "protocol": "tcp",
       "to_port": 80
+    },
+    {
+      "cidr_blocks": [
+        "0.0.0.0/0"
+      ],
+      "from_port": 443,
+      "protocol": "tcp",
+      "to_port": 443
     }
   ],
   "name": "allow_http"
 }
````````

</details>
