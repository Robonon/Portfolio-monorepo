terraform {
  required_providers {
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

resource "helm_release" "kgateway-crds" {
  name = "kgateway-crds"
  chart = "${path.module}/charts/kgateway-crds-v2.1.1.tgz"
}

# resource "helm_release" "kgateway" {
#   depends_on = [ helm_release.kgateway-crds ]
#   name = "kgateway"
#   chart = "${path.module}/charts/kgateway-v2.1.1.tgz"
#   namespace  = "kgateway-system"
#   create_namespace = true
#   values = [file("${path.module}/values/kgateway-values.yaml")]
#   force_update = true
#   cleanup_on_fail = true
#   atomic = true
#   replace = true
#   dependency_update = true
#   wait = true
# }

resource "helm_release" "argocd" {
  name       = "argocd"
  recreate_pods = true
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "9.1.4"
  namespace  = "argocd"
  create_namespace = true
  values     = [file("${path.module}/values/argocd-values.yaml")]
}

resource "helm_release" "argocd_resources" {
  # depends_on = [ helm_release.argocd ]
  recreate_pods = true
  name       = "argocd-resources"
  chart      = "${path.module}/argocd-resources/"
  namespace  = "argocd"
  create_namespace = true
  values = [yamlencode(local.argocd_resources_values)]
}

resource "tls_private_key" "tls_private_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "tls_cert" {
  private_key_pem = tls_private_key.tls_private_key.private_key_pem
  subject {
    common_name  = "argocd.localhost"
  }
  validity_period_hours = 8760
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}