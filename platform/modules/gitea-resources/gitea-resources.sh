#!/usr/bin/env bash

# Config
GITEA_URL="http://localhost:3000"
API_TOKEN="your_gitea_token"
ORG="company"
TEAM="platform"
REPO="platform-ops"
DEPLOY_KEY_TITLE="platform-ops-deploy-key"
DEPLOY_KEY_PATH="./platform-ops-deploy-key.pub" # Path to your public key

# 1. Create org (if not exists)
curl -s -X POST "$GITEA_URL/api/v1/orgs" \
  -H "Authorization: token $API_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$ORG\"}"

# 2. Create team (if not exists)
curl -s -X POST "$GITEA_URL/api/v1/orgs/$ORG/teams" \
  -H "Authorization: token $API_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"$TEAM\",\"description\":\"Platform Team\",\"permission\":\"owner\"}"

# 3. Create repo (if not exists)
curl -s -X POST "$GITEA_URL/api/v1/org/$ORG/repos" \
  -H "Authorization: token $API_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"$REPO\",\"private\":false,\"auto_init\":true}"

# # 4. Add deploy key
# PUBKEY=$(cat "$DEPLOY_KEY_PATH")
# curl -s -X POST "$GITEA_URL/api/v1/repos/$ORG/$REPO/keys" \
#   -H "Authorization: token $API_TOKEN" \
#   -H "Content-Type: application/json" \
#   -d "{\"title\":\"$DEPLOY_KEY_TITLE\",\"key\":\"$PUBKEY\",\"read_only\":false}"