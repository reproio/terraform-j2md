test data
---

`terraform plan` 実行時のログとそれに対応する JSON を設置する予定
各ディレクトリで下記を実行すると出力のし直しができる.

```
terraform init
terraform plan -no-color -out /tmp/tfplan | tee plan.txt && terraform show -json /tmp/tfplan | tee show.json
```
