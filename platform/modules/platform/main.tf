terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
      version = "0.7.0"
    }
    argocd = {
      source = "argoproj-labs/argocd"
      version = "7.12.0"
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

# Kustomize application
resource "argocd_application" "kustomize" {
  metadata {
    name      = "kustomize-app"
    namespace = "argocd"
    labels = {
      test = "true"
    }
  }

  cascade = false # disable cascading deletion
  wait    = true

  spec {
    project = "myproject"

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "foo"
    }

    source {
      repo_url        = "https://github.com/kubernetes-sigs/kustomize"
      path            = "examples/helloWorld"
      target_revision = "master"
      kustomize {
        name_prefix = "foo-"
        name_suffix = "-bar"
        images      = ["hashicorp/terraform:light"]
        common_labels = {
          "this.is.a.common" = "la-bel"
          "another.io/one"   = "true"
        }
      }
    }

    sync_policy {
      automated {
        prune       = true
        self_heal   = true
        allow_empty = true
      }
      # Only available from ArgoCD 1.5.0 onwards
      sync_options = ["Validate=false"]
      retry {
        limit = "5"
        backoff {
          duration     = "30s"
          max_duration = "2m"
          factor       = "2"
        }
      }
    }

    ignore_difference {
      group         = "apps"
      kind          = "Deployment"
      json_pointers = ["/spec/replicas"]
    }

    ignore_difference {
      group = "apps"
      kind  = "StatefulSet"
      name  = "someStatefulSet"
      json_pointers = [
        "/spec/replicas",
        "/spec/template/spec/metadata/labels/bar",
      ]
      # Only available from ArgoCD 2.1.0 onwards
      jq_path_expressions = [
        ".spec.replicas",
        ".spec.template.spec.metadata.labels.bar",
      ]
    }
  }
}

# Helm application
resource "argocd_application" "helm" {
  metadata {
    name      = "helm-app"
    namespace = "argocd"
    labels = {
      test = "true"
    }
  }

  spec {
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "default"
    }

    source {
      repo_url        = "https://some.chart.repo.io"
      chart           = "mychart"
      target_revision = "1.2.3"
      helm {
        release_name = "testing"
        parameter {
          name  = "image.tag"
          value = "1.2.3"
        }
        parameter {
          name  = "someotherparameter"
          value = "true"
        }
        value_files = ["values-test.yml"]
        values = yamlencode({
          someparameter = {
            enabled   = true
            someArray = ["foo", "bar"]
          }
        })
      }
    }
  }
}

resource "argocd_application" "platform_apps" {
  metadata {
    name      = "platform-apps"
    namespace = "argocd"
    labels = {
      platform = "platform"
    }
  }

  spec {
    project = "platform-project"

    source {
        repo_url = "https://github.com/Robonon/Portfolio-monorepo.git"
        chart = "mychart"
        path = "/apps/platform-apps"
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "default"
    }
  }
}