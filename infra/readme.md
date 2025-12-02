├── infra/                        # Terraform (cloud + cluster)
│   ├── global/                   # org-level resources (DNS, accounts, S3 buckets)
│   ├── envs/
│   │   ├── dev/
│   │   ├── staging/
│   │   └── prod/
│   └── modules/                  # Terraform reusable modules



## Start k3d cluster ( doesn't seem to work on windows right now )
requires:
- k3d

run:
- k3d cluster create --config k3d.yaml
- k3d cluster delete <cluster> (to delete cluster)

## Start kind cluster
requires:
- kind

run:
- create-kind-cluster.sh

## TODO
- install argocd in cluster