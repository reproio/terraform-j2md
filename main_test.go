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
	// # null_resource.foo will be created
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
	// # null_resource.foo will be destroyed
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
	// # env_variable.test1 will be updated in-place
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
	// # random_id.test will be replaced
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
	// # env_variable.test2 will be updated in-place
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
	// # env_variable.test3 will be destroyed
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
	// # env_variable.test5 will be created
	// @@ -1 +1,3 @@
	// -null
	// +{
	// +  "name": "test5"
	// +}
	// ```
	//
	// ```diff
	// # random_id.test4 will be replaced
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

func ExampleAwsMixed() {
	RunTestMain("test_data/aws_mixed/show.json")
	// Output:
	// ### 3 to add, 1 to change, 2 to destroy.
	// - add
	// 	- aws_route_table.public-route
	// 	- aws_route_table_association.puclic-a
	// - change
	// 	- aws_subnet.public-a
	// - destroy
	// 	- aws_instance.test
	// - replace
	// 	- aws_security_group.admin
	// <details><summary>Change details (Click me)</summary>
	//
	// ```diff
	// # aws_instance.test will be destroyed
	// @@ -1,92 +1 @@
	// -{
	// -  "ami": "ami-cbf90ecb",
	// -  "arn": "arn:aws:ec2:ap-northeast-1:999999999999:instance/i-0ecc384fa6f8d0623",
	// -  "associate_public_ip_address": false,
	// -  "availability_zone": "ap-northeast-1a",
	// -  "capacity_reservation_specification": [
	// -    {
	// -      "capacity_reservation_preference": "open",
	// -      "capacity_reservation_target": []
	// -    }
	// -  ],
	// -  "cpu_core_count": 1,
	// -  "cpu_threads_per_core": 1,
	// -  "credit_specification": [
	// -    {
	// -      "cpu_credits": "standard"
	// -    }
	// -  ],
	// -  "disable_api_termination": false,
	// -  "ebs_block_device": [],
	// -  "ebs_optimized": false,
	// -  "enclave_options": [
	// -    {
	// -      "enabled": false
	// -    }
	// -  ],
	// -  "ephemeral_block_device": [],
	// -  "get_password_data": false,
	// -  "hibernation": false,
	// -  "host_id": null,
	// -  "iam_instance_profile": "",
	// -  "id": "i-0ecc384fa6f8d0623",
	// -  "instance_initiated_shutdown_behavior": "stop",
	// -  "instance_state": "running",
	// -  "instance_type": "t2.micro",
	// -  "ipv6_address_count": 0,
	// -  "ipv6_addresses": [],
	// -  "key_name": "id_rsa_ec2",
	// -  "launch_template": [],
	// -  "metadata_options": [
	// -    {
	// -      "http_endpoint": "enabled",
	// -      "http_put_response_hop_limit": 1,
	// -      "http_tokens": "optional",
	// -      "instance_metadata_tags": "disabled"
	// -    }
	// -  ],
	// -  "monitoring": false,
	// -  "network_interface": [],
	// -  "outpost_arn": "",
	// -  "password_data": "",
	// -  "placement_group": "",
	// -  "placement_partition_number": null,
	// -  "primary_network_interface_id": "eni-081e509528cb47cc0",
	// -  "private_dns": "ip-10-1-1-11.ap-northeast-1.compute.internal",
	// -  "private_ip": "10.1.1.11",
	// -  "public_dns": "",
	// -  "public_ip": "",
	// -  "root_block_device": [
	// -    {
	// -      "delete_on_termination": true,
	// -      "device_name": "/dev/xvda",
	// -      "encrypted": false,
	// -      "iops": 100,
	// -      "kms_key_id": "",
	// -      "tags": {},
	// -      "throughput": 0,
	// -      "volume_id": "vol-072b863083c3ea911",
	// -      "volume_size": 8,
	// -      "volume_type": "gp2"
	// -    }
	// -  ],
	// -  "secondary_private_ips": [],
	// -  "security_groups": [],
	// -  "source_dest_check": true,
	// -  "subnet_id": "subnet-0342dca4d2a611266",
	// -  "tags": {
	// -    "Name": "test_ec2"
	// -  },
	// -  "tags_all": {
	// -    "Name": "test_ec2"
	// -  },
	// -  "tenancy": "default",
	// -  "timeouts": null,
	// -  "user_data": null,
	// -  "user_data_base64": null,
	// -  "user_data_replace_on_change": false,
	// -  "volume_tags": null,
	// -  "vpc_security_group_ids": [
	// -    "sg-05bf69021f9e927aa"
	// -  ]
	// -}
	// +null
	// ```
	//
	// ```diff
	// # aws_route_table.public-route will be created
	// @@ -1 +1,22 @@
	// -null
	// +{
	// +  "route": [
	// +    {
	// +      "carrier_gateway_id": "",
	// +      "cidr_block": "0.0.0.0/0",
	// +      "destination_prefix_list_id": "",
	// +      "egress_only_gateway_id": "",
	// +      "gateway_id": "igw-0edc99b3ee0ed84ad",
	// +      "instance_id": "",
	// +      "ipv6_cidr_block": "",
	// +      "local_gateway_id": "",
	// +      "nat_gateway_id": "",
	// +      "network_interface_id": "",
	// +      "transit_gateway_id": "",
	// +      "vpc_endpoint_id": "",
	// +      "vpc_peering_connection_id": ""
	// +    }
	// +  ],
	// +  "tags": null,
	// +  "timeouts": null,
	// +  "vpc_id": "vpc-0c08ee65bf93a360f"
	// +}
	// ```
	//
	// ```diff
	// # aws_route_table_association.puclic-a will be created
	// @@ -1 +1,4 @@
	// -null
	// +{
	// +  "gateway_id": null,
	// +  "subnet_id": "subnet-0342dca4d2a611266"
	// +}
	// ```
	//
	// ```diff
	// # aws_security_group.admin will be replaced
	// @@ -1,6 +1,5 @@
	//  {
	// -  "arn": "arn:aws:ec2:ap-northeast-1:999999999999:security-group/sg-05bf69021f9e927aa",
	// -  "description": "test",
	// +  "description": "description",
	//    "egress": [
	//      {
	//        "cidr_blocks": [
	// @@ -16,7 +15,6 @@
	//        "to_port": 0
	//      }
	//    ],
	// -  "id": "sg-05bf69021f9e927aa",
	//    "ingress": [
	//      {
	//        "cidr_blocks": [
	// @@ -33,11 +31,8 @@
	//      }
	//    ],
	//    "name": "admin",
	// -  "name_prefix": "",
	// -  "owner_id": "999999999999",
	//    "revoke_rules_on_delete": false,
	// -  "tags": {},
	// -  "tags_all": {},
	// +  "tags": null,
	//    "timeouts": null,
	//    "vpc_id": "vpc-0c08ee65bf93a360f"
	//  }
	// ```
	//
	// ```diff
	// # aws_subnet.public-a will be updated in-place
	// @@ -18,10 +18,10 @@
	//    "owner_id": "999999999999",
	//    "private_dns_hostname_type_on_launch": "ip-name",
	//    "tags": {
	// -    "Name": "test_subnet"
	// +    "Name": "test_subnet1"
	//    },
	//    "tags_all": {
	// -    "Name": "test_subnet"
	// +    "Name": "test_subnet1"
	//    },
	//    "timeouts": null,
	//    "vpc_id": "vpc-0c08ee65bf93a360f"
	// ```
	//
	// </details>
}
