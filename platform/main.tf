provider "helm" {
  kubernetes = {
  config_path = "~/.kube/config"
  }
}

resource "helm_release" "kgateway-crds" {
  name = "kgateway-crds"
  chart      = "${path.module}/charts/kgateway-crds-v2.1.1.tgz"
  namespace  = "kgateway-system"
  create_namespace = true
  # values = [  ]
}

resource "helm_release" "kgateway" {
  depends_on = [ helm_release.kgateway-crds ]
  name = "kgateway"
  chart = "${path.module}/charts/kgateway-v2.1.1.tgz"
  # namespace  = "kgateway-system"
  create_namespace = true
  values = [file("${path.module}/values/kgateway-values.yaml")]
  force_update = true
  cleanup_on_fail = true
  atomic = true
  replace = true
  dependency_update = true
  wait = true
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