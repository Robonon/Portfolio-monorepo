platform/                        # Platform-level config, charts, and Terraform modules
├── add-hosts.ps1                # Script to add host entries (Windows)
├── add-hosts.sh                 # Script to add host entries (Linux/macOS)
├── create-module.sh             # Script to scaffold new modules
├── readme.md                    # This file
├── terraform.tf                 # Terraform config (entrypoint)
├── modules/                     # Custom Terraform modules and charts
│   ├── gitea-resources/         # Gitea provisioning module
│   ├── helm-charts/             # Helm chart packaging, ArgoCD resources, and chart templates

## Terraform Commands

Run these from the `platform` directory:
- `terraform init`      # Initialize Terraform
- `terraform plan`      # Preview changes
- `terraform apply`     # Apply changes
- `terraform destroy`   # Destroy resources

## Services Running in Cluster

- To access services externally, port-forward the gateway (namespace: `gateway`):

- **Argo CD UI:** [http://localhost:8080](http://localhost:8080)
    - Login: `admin`
    - Password: `admin`

- **Gitea Server:** [http://localhost:3000](http://localhost:3000)
    - Login: `gitea_admin`
    - Password: `gitea_admin`

- **Platform Apps:** [http://platform-apps:8000](http://platform-apps:8000)

## Notes

- Chart sources and values are under `modules/helm-charts/` and `modules/helm-charts/values/`.
- Service names and namespaces may differ from older docs; use `kubectl get svc -A` to confirm.
- If you see CrashLoopBackOff for gateway pods, check logs and config:
    ```powershell
    kubectl -n gateway logs <pod> -c kgateway-proxy
    kubectl -n gateway get configmap gateway -o yaml
    kubectl -n gateway describe pod <pod>
    ```
- Cluster may use IPv6 PodIPs (e.g., `fd00:...`). Ensure gateway/Envoy listeners bind to `::` for health probes.

## Where to look next

- `modules/helm-charts/gateway/` — gateway chart templates and values
- `modules/gitea-resources/` — Terraform for Gitea provisioning