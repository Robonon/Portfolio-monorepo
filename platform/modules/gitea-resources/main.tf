terraform {
    required_providers {
      gitea = {
            source = "go-gitea/gitea"
            version = "0.7.0"
        }
    }
}
data "gitea_repo" "platform_ops_repo" {
  username = "company"
  name  = "platform-ops"
}

resource "tls_private_key" "platform_ops_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "gitea_repository_key" "platform_ops_key" {
  title = "platform-ops-deploy-key"
  key = tls_private_key.platform_ops_key.public_key_openssh
  read_only = false
  repository = data.gitea_repo.platform_ops_repo.id
}