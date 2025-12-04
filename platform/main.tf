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
  }
}
provider "argocd" {
  # Configuration options
}
provider "helm" {
  kubernetes = {
  config_path = "~/.kube/config"
  }
}
provider "kubernetes" {
  config_path = "~/.kube/config"
}


resource "helm_release" "kgateway-crds" {
  name = "kgateway-crds"
  chart      = "${path.module}/charts/kgateway-crds-v2.1.1.tgz"
  # namespace  = "kgateway-system"
  # values = [  ]
}

resource "helm_release" "kgateway" {
  depends_on = [ helm_release.kgateway-crds ]
  name = "kgateway"
  chart = "${path.module}/charts/kgateway-v2.1.1.tgz"
  namespace  = "kgateway-system"
  values = [file("${path.module}/values/kgateway-values.yaml")]
  force_update = true
  cleanup_on_fail = true
  atomic = true
  replace = true
  dependency_update = true
  wait = true
}

resource "kubernetes_manifest" "gateway" {
  manifest = yamldecode(file("${path.module}/templates/gateway.yaml"))
  depends_on = [ helm_release.kgateway ]
}

resource "kubernetes_manifest" "argocd_route" {
  manifest = yamldecode(file("${path.module}/templates/argocd-route.yaml"))
  depends_on = [ kubernetes_manifest.gateway ]
}

resource "kubernetes_manifest" "gitea_route" {
  manifest = yamldecode(file("${path.module}/templates/gitea-route.yaml"))
  depends_on = [ kubernetes_manifest.gateway ]
}


resource "helm_release" "gitea_server" {
  name       = "gitea-server"
  repository = "https://dl.gitea.com/charts/"
  chart      = "gitea"
  version    = "12.4.0"
  namespace  = "gitea-server"
  create_namespace = true
  values     = [file("${path.module}/values/gitea-server-values.yaml")]
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "9.1.4"
  namespace  = "argocd"
  create_namespace = true
  values     = [file("${path.module}/values/argocd-values.yaml")]
}

