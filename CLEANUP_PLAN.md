# ScoutTalent Platform - Cleanup Plan

## 🔍 Audit Results

### Current Structure Analysis

```
workspace/
├── .github/workflows/          ✅ KEEP - CI/CD pipelines
├── app/frontend/               ❌ DELETE - Empty except node_modules
├── backend/                    ✅ KEEP - All Go services
│   ├── services/              ✅ KEEP (10 services total)
│   │   ├── ai-moderation-worker/    ✅ NEW - Keep
│   │   ├── ai-scoring-worker/       ⚠️  REVIEW - Not documented
│   │   ├── auth-service/            ✅ KEEP
│   │   ├── discovery-service/       ✅ NEW - Keep
│   │   ├── engagement-service/      ⚠️  REVIEW - Not documented
│   │   ├── highlight-generator/     ⚠️  REVIEW - Not documented
│   │   ├── media-service/           ✅ KEEP
│   │   ├── notification-service/    ⚠️  REVIEW - Not documented
│   │   ├── payment-service/         ⚠️  REVIEW - Not documented
│   │   └── profile-service/         ✅ KEEP
│   ├── pkg/                   ✅ KEEP - Shared packages
│   ├── docs/                  ✅ KEEP - Documentation
│   ├── scripts/               ✅ KEEP - Test scripts
│   ├── docker-compose.yml     ✅ KEEP
│   ├── Makefile               ✅ KEEP
│   └── *.md files             ✅ KEEP
├── frontend/                   ✅ KEEP - Nuxt 3 app
├── helm/                       ⚠️  REVIEW - Kubernetes deployment
├── k8s/                        ⚠️  REVIEW - Kubernetes manifests
├── terraform/                  ⚠️  REVIEW - Infrastructure as code
├── tests/                      ❌ DELETE - Empty directory
├── web/                        ❌ DELETE - Old frontend (redundant)
├── README.md                   ✅ KEEP
└── PROJECT_COMPLETION_SUMMARY.md ✅ KEEP
```

## 🗑️ Items to Delete

### 1. **app/frontend/** - REDUNDANT
- Only contains empty node_modules
- Real frontend is in `/frontend/`
- **Action**: Delete entire directory

### 2. **web/** - OLD FRONTEND
- Contains old React/Vue frontend
- Replaced by new Nuxt 3 frontend in `/frontend/`
- **Action**: Delete entire directory OR move to archive

### 3. **tests/** - EMPTY
- Empty directory
- Tests are in individual service directories
- **Action**: Delete

## ⚠️ Items to Review

### Infrastructure Directories
These contain Kubernetes and Terraform configs that may be needed for deployment:

1. **helm/** - Helm charts for Kubernetes
2. **k8s/** - Raw Kubernetes manifests
3. **terraform/** - Infrastructure as Code

**Recommendation**: 
- If deploying to Kubernetes: KEEP and move to `backend/deployment/`
- If not using Kubernetes yet: MOVE to archive or separate deployment repo

### Undocumented Services
These services exist but aren't in the main documentation:

1. **ai-scoring-worker** - Purpose unclear
2. **engagement-service** - Purpose unclear
3. **highlight-generator** - Purpose unclear
4. **notification-service** - Purpose unclear
5. **payment-service** - Purpose unclear

**Recommendation**: 
- Review each service's README
- Determine if they're:
  - a) Active and needed → Document in main README
  - b) Incomplete → Move to `backend/services-wip/`
  - c) Obsolete → Delete

## ✅ Cleanup Actions

### Phase 1: Safe Deletions (No Risk)
```bash
# Delete empty/redundant directories
rm -rf app/
rm -rf tests/
```

### Phase 2: Archive Old Frontend (Low Risk)
```bash
# Option A: Delete if confirmed not needed
rm -rf web/

# Option B: Archive for reference
mkdir -p archive/
mv web/ archive/web-old-frontend/
```

### Phase 3: Organize Infrastructure (Medium Risk)
```bash
# Move deployment configs to backend
mkdir -p backend/deployment/
mv helm/ backend/deployment/helm/
mv k8s/ backend/deployment/k8s/
mv terraform/ backend/deployment/terraform/
```

### Phase 4: Review Services (Requires Decision)
```bash
# Create WIP directory for incomplete services
mkdir -p backend/services-wip/

# Move undocumented services (if not needed immediately)
# Example:
# mv backend/services/ai-scoring-worker backend/services-wip/
```

## 📊 Expected Results After Cleanup

### Clean Structure:
```
scouttalent-platform/
├── .github/workflows/          # CI/CD
├── backend/
│   ├── services/              # 5 core services
│   │   ├── auth-service/
│   │   ├── profile-service/
│   │   ├── media-service/
│   │   ├── ai-moderation-worker/
│   │   └── discovery-service/
│   ├── services-wip/          # Work in progress (optional)
│   ├── deployment/            # K8s, Helm, Terraform
│   ├── pkg/
│   ├── docs/
│   └── scripts/
├── frontend/                   # Nuxt 3 app
├── archive/                    # Old code (optional)
├── README.md
└── PROJECT_COMPLETION_SUMMARY.md
```

### Benefits:
- ✅ Clear separation of concerns
- ✅ No redundant directories
- ✅ Easier navigation
- ✅ Smaller repository size
- ✅ Clearer documentation

## 🎯 Recommended Immediate Actions

1. **Delete `app/` directory** - Confirmed redundant
2. **Delete `tests/` directory** - Empty
3. **Review `web/` directory** - Likely old frontend, can delete
4. **Document or archive** undocumented services
5. **Organize deployment configs** into `backend/deployment/`

## 📝 Questions to Answer

1. **Are you using Kubernetes?** 
   - Yes → Keep helm/, k8s/, terraform/ but organize them
   - No → Archive them for future use

2. **What are these services for?**
   - ai-scoring-worker
   - engagement-service
   - highlight-generator
   - notification-service
   - payment-service
   
3. **Is the `web/` directory needed?**
   - Likely old frontend that can be deleted

Would you like me to proceed with the cleanup?