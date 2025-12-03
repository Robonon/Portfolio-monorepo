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

resource "kubernetes_manifest" "routes" {
  manifest = yamldecode(file("${path.module}/templates/routes.yaml"))
  depends_on = [ kubernetes_manifest.gateway ]
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