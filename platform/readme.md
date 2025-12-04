├── platform/                     # Everything running *on* the cluster (Ops layer)
│   ├── charts                    # Charts that for some reason needs to be locally available
│   ├── Values/                   # Values files for the charts defined in main.tf   
│   │   ├── argocd-values.yaml 
│   │   ├── kgateway-values.yaml
|   ├── main.tf                   # Defines all cluster addons/tools in the form of helm charts


## Terraform

commands:
- terraform init
- terraform plan
- terraform apply
- terraform destory

## Services running in cluster

- kubectl port-forward svc/gateway 8080:8000 -n kgateway-system

- [argocd ui](http://argocd.localhost:8080)
    - login: admin, get password with 'argocd admin initial-password -n argocd'
- [gitea server](http://gitea.localhost:8080)
    - login: gitea_admin, r8sA8CPHD9!bt6d