# modules/helm-charts/variables.tf
variable "platform_ops_repo_deploy_key" {
  type      = string
  sensitive = true
}

variable "platform_project_name" {
  default = "platform-project"
}

variable "gateway_project_name" {
  default = "gateway-project"
}

variable "argocd_project_name" {
  default = "argocd-project"
}

locals {
  argocd_resources_values = {
    // Projects
    projects = {
      platform = {
        name      = var.platform_project_name
        namespace = "platform-apps"
      }
      gateway = {
        name      = var.gateway_project_name
        namespace = "gateway"
      }
      argocd = {
        name      = var.argocd_project_name
        namespace = "argocd"
      }
    }
    // Repositories
    repositories = {
      platform-apps-chart = {
        name    = "platform-apps-chart"
        url     = "oci://host.docker.internal:5000/platform-apps"
        project = var.platform_project_name
        type    = "oci"
      }
      gateway-chart = {
        name    = "gateway-chart"
        url     = "oci://host.docker.internal:5000/gateway"
        project = var.gateway_project_name
        type    = "oci"
      }
      tenant-image = {
        name    = "tenant-image"
        url     = "oci://host.docker.internal:5000/tenant"
        project = var.platform_project_name
        type    = "oci"
      }
    }
    // Applications
    applications = {
      argocd = {
        name           = "argocd"
        project        = var.argocd_project_name
        chart          = "argo-cd"
        repoURL        = "https://argoproj.github.io/argo-helm"
        targetRevision = "9.1.4"
        namespace      = "argocd"
      }
      gateway = {
        name           = "gateway"
        project        = var.gateway_project_name
        chart          = "gateway"
        repoURL        = "oci://host.docker.internal:5000/gateway"
        targetRevision = "0.1.1"
        namespace      = "gateway"
        config = {
          "tls.cert" = tls_self_signed_cert.tls_cert.cert_pem
          "tls.key"  = tls_private_key.tls_private_key.private_key_pem
        }
      }
      platform = {
        name           = "platform-apps"
        project        = var.platform_project_name
        chart          = "platform-apps"
        repoURL        = "oci://host.docker.internal:5000/platform-apps"
        targetRevision = "0.1.1"
        namespace      = "platform-apps"
      }
    }
  }
}