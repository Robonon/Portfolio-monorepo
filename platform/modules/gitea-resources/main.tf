terraform {
    required_providers {
      gitea = {
            source = "go-gitea/gitea"
            version = "0.7.0"
        }
    }
}

data "gitea_org" "company_data" {
  name = "company"
}

resource "gitea_org" "company" {
  name = data.gitea_org.company_data.name
}

resource "gitea_team" "platform_team" {
#   depends_on = [ gitea_org.company ]
  name         = "platform"
  organisation = data.gitea_org.company_data.name
  description  = "Platform Team"
  permission   = "owner"
}

resource "gitea_repository" "platform_ops" {
  depends_on = [ gitea_org.company ]
  username     = data.gitea_org.company_data.name
  name         = "platform-ops"
  private      = false
  issue_labels = "Default"
  lifecycle {
    prevent_destroy = false
  }
}