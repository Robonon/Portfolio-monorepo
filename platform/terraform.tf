terraform {
  required_providers {
    argocd = {
      source = "argoproj-labs/argocd"
      version = "7.12.0"
    }
    helm = {
      source = "hashicorp/helm"
      version = "3.1.1"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.38.0"
    }
    gitea = {
      source = "go-gitea/gitea"
      version = "0.7.0"
    }
  }
}

// Providers, global for all modules
provider "helm" {
  kubernetes = {
    config_path = "~/.kube/config"
  }
}
provider "kubernetes" {
  config_path = "~/.kube/config"
}
provider "argocd" {
  # Configuration options
}

provider "gitea" {
  base_url = "http://localhost:3000"
  username = "gitea_admin"
  password = "gitea_admin"
}

// cluster resources and addons
module "infra" {
  source = "./modules/infra"
}

// Team resources module per team
module "platform" {
  source = "./modules/platform"
}