terraform plan -var "do_token=${DO_PAT}" -var "pvt_key=$HOME/.ssh/id_rsa"

terraform apply -var "do_token=${DO_PAT}" -var "pvt_key=$HOME/.ssh/id_rsa"

terraform show

terraform plan -destroy -out=terraform.tfplan   -var "do_token=${DO_PAT}"   -var "pvt_key=$HOME/.ssh/id_rsa"

terraform apply "terraform.tfplan"

https://medium.com/@orestovyevhen/set-up-infrastructure-in-hetzner-cloud-using-terraform-ce85491e92d
https://maddevs.io/blog/terraform-hetzner/