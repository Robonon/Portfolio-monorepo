terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
      version = "0.7.0"
    }
  }
}

data "gitea_org" "company_data" {
  name = "Company"
}

resource "gitea_org" "company" {
  name = data.gitea_org.company_data.name
}

resource "gitea_repository" "platform_ops" {
#   depends_on = [ gitea_org.company ]
  username     = "gitea_admin"
  name         = "platform-ops"
  private      = false
  issue_labels = "Default"
  lifecycle {
    prevent_destroy = true // Critical to prevent accidental deletion
  }
}

resource "gitea_team" "platform_team" {
#   depends_on = [ gitea_org.company ]
  name         = "Platform"
  organisation = data.gitea_org.company_data.name
  description  = "Platform Team"
  permission   = "owner"
}