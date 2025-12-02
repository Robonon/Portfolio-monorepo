├── platform/                     # Everything running *on* the cluster (Ops layer)
│   ├── helm/                     # Helm charts for platform components
│   ├── kustomize/                # (optional) if you prefer Kustomize for operators
│   ├── addons/
│   │   ├── ingress-nginx/
│   │   ├── cert-manager/
│   │   ├── metrics-server/
│   │   ├── external-dns/
│   │   ├── argocd/
│   │   └── logging/              # Loki, Grafana, etc
│   └── values/                   # environment-specific overrides for platform charts