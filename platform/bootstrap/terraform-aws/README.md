# terraform files

This folder contains the `*.tf` files read by Terraform.
Place your `terraform.tfvars` file with your custom values, don't update the `variables.tf` file.
Example content:

```
aws_name = "tgm-ikt"
bootstrap_key_name = "tgm-us-west-2"
aws_profile = "default"
infrakit_config_base_url = "https://raw.githubusercontent.com/appcelerator/amp/branchname/bootstrap"
```
