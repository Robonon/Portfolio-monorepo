terraform {
    required_providers {
        gitea = {
            source = "go-gitea/gitea"
            version = "0.7.0"
        }
        # argocd = {
        #     source = "argoproj-labs/argocd"
        #     version = "7.12.0"
        # }
        helm = {
            source = "hashicorp/helm"
            version = "3.1.1"
        }
        kubernetes = {
            source = "hashicorp/kubernetes"
            version = "3.0.0"
        }
    }
}

provider "helm" {
    kubernetes = {
        config_path = "~/.kube/config"
    }
}

provider "kubernetes" {
    config_path = "~/.kube/config"
}

# provider "argocd" {
#     server_addr = "argocd.localhost:8443"
#     username    = "admin"
#     password    = "admin"
#     insecure    = true
# }

provider "gitea" {
    base_url = "http://localhost:3000"
    username = "gitea_admin"
    password = "gitea_admin"
}

# module "gitea-resources" {
#     source = "./modules/gitea-resources" 
# }

module "helm-charts" {
    source = "./modules/helm-charts"
}
