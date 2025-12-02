├── platform/                     # Everything running *on* the cluster (Ops layer)
│   ├── charts                    # Charts that for some reason needs to be locally available
│   ├── Values/                   # Values files for the charts defined in main.tf   
│   │   ├── argocd-values.yaml 
│   │   ├── kgateway-values.yaml
|   ├── main.tf                   # Defines all cluster addons/tools in the form of helm charts