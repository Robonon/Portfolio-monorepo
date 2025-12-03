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
Each folder contains it's own readme, explaining details about it's content

repo/
├── apps/                         # Bussines applications
├── docs/                         # docs for everything in the repo
├── gitops/                       # App of apps repo used by argocd to deploy all apps or multitenant solutions
├── infra/                        # Terraform (cloud + cluster)
├── platform/                     # Everything running *on* the cluster (Ops layer)

## Commit messages

Ticket id (in the future)
- [Some-id]

Action
- '[Add]'
- '[Update]'

Full example

git commit -m "[Some-Id-123][Add]: Add some feature or thing"