# Portfolio-monorepo

## Project ideas:

- Ops ui to manage cluster and applications
    - add new cluster
    - remove cluster
    - add or move tenant to cluster
    - remove tenant from cluster
    - visualize tenant/cluster logging
    - visualize tenant/cluster metrics

- Git server
    - Code repo hosting
    - Self-hosted runners for pipelines

- LLM and MCP servers
    - Ollama
    - Client interface
    - MCP servers
        - File manipulation
        - Database query
        - Metrics and log analysis
        - Git integration
        - Project planning

- Various business applications

- Gitops repo
    - Derive and provision resource needs from apps (event bus, db, other external resources)

- IaC
    - Local Cluster (kind, k3d)

- Config as Code
    - Sync app config schemas with infra/ops tools, to make explicit what is configurable
    

## Cluster 

- package apps with helm
- define resources with terraform

Ingress -> Ops ui       ->
        -> mcp client   -> api gateway -> queue/service endpoints -> resource endpoints

resources: 

## Repo structure
repo/
├── infra/                        # Terraform (cloud + cluster)
│   ├── global/                   # org-level resources (DNS, accounts, S3 buckets)
│   ├── envs/
│   │   ├── dev/
│   │   ├── staging/
│   │   └── prod/
│   └── modules/                  # Terraform reusable modules
│
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
│
├── apps/                         # Applications developed by you
│   ├── service-a/
│   │   ├── src/
│   │   ├── Dockerfile
│   │   ├── charts/               # Helm chart for this service
│   │   └── values/
│   ├── service-b/
│   └── shared/                   # shared libs, templates, configs
│
├── environments/                 # GitOps definitions (if using ArgoCD/Flux)
│   ├── dev/
│   │   ├── apps/
│   │   └── platform/
│   ├── staging/
│   └── prod/
│
├── ci/                           # CI/CD workflows (GitHub Actions, GitLab, etc)
└── scripts/                      # helper scripts (bootstrap, linting, local dev)


## Commit messages

Ticket id (in the future)
- [Some-id]

Action
- '[Add]'
- '[Update]'

Full example

git commit -m "[Some-Id-123][Add]: Add some feature or thing"