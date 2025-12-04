
platform/                        # Everything running *on* the cluster (Ops layer)
├── add-hosts.ps1                 # Script to add host entries (Windows)
├── add-hosts.sh                  # Script to add host entries (Linux/macOS)
├── create-module.sh              # Script to scaffold new modules
├── readme.md                     # This file
├── terraform.tf                  # Terraform config (legacy or entrypoint)
├── modules/                      # Custom Terraform modules
│   ├── platform/                 # Platform module
│   │   ├── main.tf
│   │   ├── outputs.tf
│   │   ├── variables.tf
│   │   ├── README.md
│   │   └── apps/
│   │       └── teams/
│   │           ├── go.mod
│   │           └── main.go
│   ├── charts/                   # Local charts (if needed)
│   └── templates/                # YAML templates for manifests
│       ├── argocd-route.yaml
│       ├── gateway.yaml
│       └── gitea-route.yaml
│   └── values/                   # Values files for charts
│       ├── argocd-values.yaml
│       ├── gitea-server-values.yaml
│       └── kgateway-values.yaml




## Terraform Commands

Run these from the `platform` directory:
- `terraform init`      # Initialize Terraform
- `terraform plan`      # Preview changes
- `terraform apply`     # Apply changes
- `terraform destroy`   # Destroy resources


## Services Running in Cluster

- To access services externally, port-forward the gateway:
    `kubectl port-forward svc/gateway 8080:8000 -n kgateway-system`

- **Argo CD UI:** [http://argocd.localhost:8080](http://argocd.localhost:8080)
    - Login: `admin`
    - Get password: `argocd admin initial-password -n argocd`

- **Gitea Server:** [http://localhost:3000](http://localhost:3000)
    - Login: `gitea_admin`
    - Password: `gitea_admin`