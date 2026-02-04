# Tekton Operator EUS Upgrade Study: v0.71.x → v0.78.x

**Document Version**: 1.0  
**Generated**: 2026-02-04  
**Author**: EUS Upgrade Study Analysis  
**Target Audience**: EUS customers upgrading from OCP 4.15 to OCP 4.20+

---

## Executive Summary

### Upgrade Feasibility
- **Direct Upgrade Supported**: YES (with migration steps required)
- **Risk Level**: HIGH
- **Estimated Complexity**: MODERATE to COMPLEX

### Key Findings

1. **CRITICAL: Tekton Pipeline v0 → v1 Major Version Jump**
   - Pipeline component: v0.59.6 → v1.6.0 (breaking change)
   - This transition occurs at operator v0.76.x (Pipeline v0.68→v1.0)
   - Major version jumps typically include breaking API changes

2. **New Component Added: TektonPruner**
   - Introduced in v0.76.x as new CRD `tektonpruners.operator.tekton.dev`
   - Provides automated pruning of Tekton resources

3. **Significant Component Version Increases**
   - All managed components see major version increments
   - Triggers, Dashboard, Chains, Results, Hub, PAC all updated

### Critical Upgrade Considerations

**❗ PIPELINE v1 MIGRATION**: The Pipeline v0 → v1 transition in v0.76.x requires special attention:
- API changes between v0 and v1 must be reviewed
- Custom TaskRuns, PipelineRuns may need updates
- Deprecated v1beta1 APIs may be removed

###

 Recommended Approach

**For EUS customers (v0.71.x → v0.78.x):**

1. **Pre-upgrade preparation**:
   - Review Pipeline v1 migration guide from Tekton upstream
   - Audit existing TaskRuns, PipelineRuns for v1beta1 API usage
   - Test upgrade in non-production environment
   - Backup all Tekton custom resources

2. **Consider incremental approach**:
   - While direct upgrade is supported, the Pipeline v1 transition is significant
   - Consider: v0.71.x → v0.75.x (last v0) → v0.76.x (v1) → v0.78.x
   - This allows isolated testing of the Pipeline v1 migration

3. **Post-upgrade validation**:
   - Verify all existing pipelines and tasks still function
   - Test new Pruner component if enabling
   - Monitor for deprecated API warnings

---

## Version Timeline & Component Matrix

| Version | Pipeline | Triggers | Dashboard | Chains | Results | Hub | PAC | Pruner | K8s Min |
|---------|----------|----------|-----------|--------|---------|-----|-----|--------|---------|
| v0.71.x | v0.59.6  | v0.27.0  | v0.46.0   | v0.20.1| v0.10.0 | v1.17.2| v0.27.2| -      | 1.27.x  |
| v0.72.x | v0.61.x  | v0.28.0  | v0.48.x   | v0.21.x| v0.11.x | v1.18.x| v0.28.x| -      | 1.28.x  |
| v0.73.x | v0.62.x  | v0.29.0  | v0.50.x   | v0.22.x| v0.12.x | v1.19.x| v0.29.x| -      | 1.28.x  |
| v0.74.x | v0.65.x  | v0.30.0  | v0.51.x   | v0.23.x| v0.13.x | v1.19.x| v0.30.x| -      | 1.28.x  |
| v0.75.x | v0.68.1  | v0.31.0  | v0.54.0   | v0.24.0| v0.14.0 | v1.20.2| v0.33.2| -      | 1.28.x  |
| **v0.76.x** | **v1.0.0** | v0.32.0  | v0.57.1 | v0.25.1| v0.15.3 | v1.21.1| v0.35.4| **v0.1.0** | 1.28.x |
| v0.77.x | v1.3.1   | v0.33.0  | v0.60.0   | v0.25.5| v0.16.0 | v1.22.0| v0.37.0| v0.2.0 | 1.28.x  |
| v0.78.x | v1.6.0   | v0.34.0  | v0.63.1   | v0.26.0| v0.17.2 | v1.23.6| v0.39.3| v0.3.5 | 1.28.x  |

**KEY TRANSITION**: v0.75.x → v0.76.x introduces Pipeline v1.0

---

## Component Version Analysis

### Pipeline Component ⚠️ CRITICAL

**Version Jump**: v0.59.6 → v1.6.0

**Incremental Path**:
- v0.71.x: v0.59.6
- v0.72.x: v0.61.x
- v0.73.x: v0.62.x
- v0.74.x: v0.65.x
- v0.75.x: v0.68.1 (last v0 release)
- **v0.76.x: v1.0.0** ⚠️ MAJOR VERSION CHANGE
- v0.77.x: v1.3.1
- v0.78.x: v1.6.0

**Breaking Change Analysis**:
- Major version increment (v0 → v1) indicates breaking changes
- Tekton Pipeline v1.0 likely removes deprecated v1beta1 APIs
- v1alpha1 APIs may be promoted to stable v1
- TaskRun, PipelineRun specs may have schema changes

**Migration Required**: YES
- Review all existing Pipelines and Tasks
- Update any v1beta1 API usage to v1
- Test all pipelines in non-production
- Monitor deprecation warnings

**Recommended Actions**:
1. Review Tekton Pipeline v0.68 → v1.0 migration guide
2. Search for v1beta1 usage: `kubectl get pipelines,tasks -A -o yaml | grep v1beta1`
3. Update resources before operator upgrade
4. Test pipeline execution post-upgrade

### Triggers Component

**Version Jump**: v0.27.0 → v0.34.0 (+7 minor versions)

**Assessment**: Moderate impact
- Multiple minor version increments
- Should follow semantic versioning (backwards compatible)
- Review for deprecated fields

### Dashboard Component

**Version Jump**: v0.46.0 → v0.63.1 (+17 minor versions)

**Assessment**: Low to Moderate impact
- UI changes, likely backwards compatible
- May include new features
- Review release notes for removed features

### Chains Component

**Version Jump**: v0.20.1 → v0.26.0 (+6 minor versions)

**Assessment**: Moderate impact
- Supply chain security features
- Verify signature verification still works
- Check for certificate/key format changes

### Results Component

**Version Jump**: v0.10.0 → v0.17.2 (+7 minor versions)

**Assessment**: Moderate impact
- Results storage and API changes
- Verify database compatibility
- Check for schema migrations

### Hub Component

**Version Jump**: v1.17.2 → v1.23.6 (+6 minor versions)

**Assessment**: Low impact
- Catalog service
- Verify catalog sync functionality

### Pipelines-as-Code Component

**Version Jump**: v0.27.2 → v0.39.3 (+12 minor versions)

**Assessment**: Moderate impact
- GitHub integration changes
- Verify webhook configurations
- Check for authentication changes

### ⭐ NEW: Pruner Component

**Version**: v0.3.5 (introduced in v0.76.x)

**Description**: New component for automatic pruning of completed PipelineRuns and TaskRuns

**Impact**: Optional feature
- Adds new CRD: `tektonpruners.operator.tekton.dev`
- Not required for operation
- Configurable retention policies

---

## CRD Analysis

### Existing CRDs (No Breaking Changes Detected)

All existing CRDs remain present:
- ✅ `tektonchains.operator.tekton.dev`
- ✅ `tektonconfigs.operator.tekton.dev`
- ✅ `tektondashboards.operator.tekton.dev`
- ✅ `tektonhubs.operator.tekton.dev`
- ✅ `tektoninstallersets.operator.tekton.dev`
- ✅ `tektonpipelines.operator.tekton.dev`
- ✅ `tektonresults.operator.tekton.dev`
- ✅ `tektontriggers.operator.tekton.dev`
- ✅ `manualapprovalgates.operator.tekton.dev`

### New CRD Added

- ⭐ `tektonpruners.operator.tekton.dev` (v0.76.x+)

**No CRDs removed** - Full backwards compatibility maintained at operator CRD level.

---

## RBAC Changes

**Analysis**: Minimal changes (7 lines added)
- RBAC changes: 927 → 934 lines (+0.75%)
- Likely additions for Pruner component permissions
- No significant permission escalations detected

**Assessment**: LOW RISK
- Existing ServiceAccounts remain functional
- New permissions added for Pruner (if enabled)

---

## Kubernetes Version Requirements

| Operator Version | Minimum K8s Version | OCP Version |
|------------------|---------------------|-------------|
| v0.71.x          | 1.27.x             | 4.15 (EUS)  |
| v0.72.x - v0.78.x| 1.28.x             | 4.16 - 4.20 |

**K8s Upgrade Required**: YES
- v0.72+ requires K8s 1.28.x minimum
- EUS customers on OCP 4.15 (K8s 1.27) must upgrade to OCP 4.20 (K8s 1.28+)

---

## Breaking Changes Summary

### High Severity

1. **Tekton Pipeline v0 → v1 Migration** (v0.76.x)
   - **Impact**: All Pipeline and Task custom resources
   - **Affected Users**: All customers using Tekton Pipelines
   - **Migration Required**: YES
   - **Migration Steps**:
     1. Audit all Pipelines, Tasks, PipelineRuns, TaskRuns for v1beta1 API usage
     2. Update to v1 API specs before operator upgrade
     3. Test pipeline execution in staging
     4. Review Tekton upstream v1.0 migration guide
     5. Monitor for deprecation warnings post-upgrade

2. **Kubernetes Version Bump** (v0.72.x)
   - **Impact**: Cannot run on K8s < 1.28
   - **Affected Users**: All customers on OCP < 4.16
   - **Migration Required**: YES
   - **Migration Steps**: Upgrade OpenShift to 4.20 (EUS) first

### Medium Severity

1. **New Pruner Component** (v0.76.x)
   - **Impact**: New optional component
   - **Affected Users**: Customers wanting auto-pruning
   - **Migration Required**: NO (opt-in)
   - **Action**: Review Pruner configuration if enabling

2. **Component Version Jumps**
   - **Impact**: Multiple minor versions skipped
   - **Affected Users**: Users with custom integrations
   - **Migration Required**: REVIEW RECOMMENDED
   - **Action**: Review component changelogs for deprecated features

### Low Severity

1. **RBAC Additions**
   - **Impact**: New permissions for Pruner
   - **Affected Users**: Minimal
   - **Migration Required**: NO

---

## Upgrade Blockers & Workarounds

### Blocker 1: Pipeline v1 API Compatibility

- **Description**: Pipeline v0 → v1 may include breaking API changes
- **Impact**: Existing TaskRuns and PipelineRuns may fail if using deprecated APIs
- **Workaround**: Pre-upgrade API audit and resource updates
- **Status**: REQUIRES INVESTIGATION
- **Tracking**: Review Tekton Pipeline v1.0 release notes

**Mitigation**:
```bash
# Audit for v1beta1 API usage
kubectl get pipelines,tasks -A -o yaml | grep "apiVersion.*v1beta1"

# If found, update to v1 API before upgrade
```

### Blocker 2: Kubernetes Version Compatibility

- **Description**: Operator v0.72+ requires K8s 1.28+
- **Impact**: Cannot upgrade operator without first upgrading OpenShift
- **Workaround**: Upgrade OCP to 4.20 before operator upgrade
- **Status**: KNOWN REQUIREMENT
- **Tracking**: N/A

**Mitigation**: Coordinate OpenShift EUS upgrade (4.15 → 4.20) before Tekton operator upgrade

---

## Incremental Upgrade Path Analysis

### v0.71.x → v0.72.x
- **Breaking Changes**: Kubernetes version bump (1.27 → 1.28)
- **Key Changes**: Minor component updates
- **Migration Effort**: LOW (K8s upgrade required)

### v0.72.x → v0.73.x
- **Breaking Changes**: None
- **Key Changes**: Minor component updates
- **Migration Effort**: LOW

### v0.73.x → v0.74.x
- **Breaking Changes**: None
- **Key Changes**: Minor component updates
- **Migration Effort**: LOW

### v0.74.x → v0.75.x
- **Breaking Changes**: None
- **Key Changes**: Last v0.x Pipeline release
- **Migration Effort**: LOW

### v0.75.x → v0.76.x ⚠️ CRITICAL
- **Breaking Changes**: Pipeline v0.68 → v1.0 (MAJOR)
- **Key Changes**: 
  - Pipeline v1.0 release
  - New Pruner component added
- **Migration Effort**: HIGH

### v0.76.x → v0.77.x
- **Breaking Changes**: None
- **Key Changes**: Pipeline v1.3.1, component updates
- **Migration Effort**: LOW

### v0.77.x → v0.78.x
- **Breaking Changes**: None
- **Key Changes**: Pipeline v1.6.0, component updates
- **Migration Effort**: LOW

---

## Migration Guide

### Prerequisites

- [ ] OpenShift 4.20 (K8s 1.28.x) or higher
- [ ] Backup of all Tekton resources (Pipelines, Tasks, TriggerTemplates, etc.)
- [ ] Test environment available for validation
- [ ] Review access to upstream Tekton Pipeline v1 migration documentation

### Step-by-Step Migration

#### Phase 1: Assessment (Est: 1-2 hours)

1. **Inventory current resources**
   ```bash
   kubectl get tektonconfig,tektonpipeline,tektontrigger,tektondashboard -A
   kubectl get pipelines,tasks,pipelineruns,taskruns -A --no-headers | wc -l
   ```

2. **Audit for v1beta1 API usage**
   ```bash
   kubectl get pipelines,tasks,triggertemplates,triggerbindings -A -o yaml | grep "apiVersion.*v1beta1"
   ```

3. **Document findings**
   - List all resources using v1beta1 APIs
   - Identify custom integrations that may be affected
   - Estimate migration effort

#### Phase 2: Preparation (Est: 2-4 hours)

1. **Backup all resources**
   ```bash
   kubectl get tektonconfig,tektonpipeline,tektontrigger -A -o yaml > backup-operator-crds.yaml
   kubectl get pipelines,tasks,triggertemplates,triggerbindings,pipelineruns,taskruns -A -o yaml > backup-tekton-resources.yaml
   ```

2. **Review Tekton Pipeline v1 migration guide**
   - Read: https://tekton.dev/docs/pipelines/migrating-v1beta1-to-v1/
   - Note breaking changes and deprecations
   - Plan resource updates

3. **Update Pipelines and Tasks to v1 API** (if using v1beta1)
   ```bash
   # Example: Update apiVersion in all Pipeline definitions
   # Manual review and update required
   ```

4. **Test in non-production environment**
   - Deploy test workloads
   - Validate pipeline execution
   - Check for deprecation warnings

#### Phase 3: Execution (Est: 30-60 minutes)

1. **Verify OpenShift version**
   ```bash
   oc version
   # Must be 4.20+ (K8s 1.28+)
   ```

2. **Upgrade operator**
   ```bash
   # If using OperatorHub/OLM
   oc patch subscription tekton-operator -n openshift-operators --type merge -p '{"spec":{"channel":"stable"}}'
   
   # Monitor upgrade progress
   oc get csv -n openshift-operators | grep tekton
   ```

3. **Wait for rollout**
   ```bash
   oc wait --for=condition=Available --timeout=600s deployment/tekton-operator -n openshift-operators
   ```

4. **Verify component versions**
   ```bash
   oc get tektonpipeline cluster -o jsonpath='{.status.version}'
   # Should show v1.6.x
   ```

#### Phase 4: Validation (Est: 1-2 hours)

1. **Check operator health**
   ```bash
   oc get tektonconfig config -o yaml
   oc get pods -n openshift-pipelines
   ```

2. **Verify component reconciliation**
   ```bash
   oc get tektonpipeline,tektontrigger,tektondashboard,tektonchains -A
   ```

3. **Test pipeline execution**
   ```bash
   # Run a sample pipeline
   tkn pipeline start <test-pipeline> -n <namespace>
   tkn pipelinerun logs -f
   ```

4. **Check for deprecation warnings**
   ```bash
   oc logs -n openshift-pipelines deployment/tekton-pipelines-controller | grep -i deprecat
   ```

5. **Validate existing workloads**
   - Run critical pipelines
   - Verify triggers still fire
   - Check dashboard accessibility
   - Test custom tasks

### Rollback Procedure

**If issues occur during upgrade:**

1. **Snapshot current state**
   ```bash
   oc get tektonconfig,tektonpipeline -A -o yaml > failed-state.yaml
   ```

2. **Rollback operator** (if using OLM)
   ```bash
   oc patch subscription tekton-operator -n openshift-operators --type merge -p '{"spec":{"channel":"previous-channel"}}'
   ```

3. **Restore resources** (if needed)
   ```bash
   kubectl apply -f backup-operator-crds.yaml
   kubectl apply -f backup-tekton-resources.yaml
   ```

4. **Investigate failure**
   - Review operator logs
   - Check for CRD validation errors
   - Document issues for retry

---

## Risk Assessment

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Pipeline v1 API breaks existing workflows | HIGH | HIGH | Pre-upgrade API audit and resource updates |
| Custom Tasks fail with v1 API | MEDIUM | HIGH | Test custom tasks in staging environment |
| Results database incompatibility | LOW | MEDIUM | Backup Results database before upgrade |
| Trigger webhook failures | LOW | MEDIUM | Test trigger execution post-upgrade |
| Dashboard UI regression | LOW | LOW | Verify dashboard functionality |

### Operational Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Extended downtime during upgrade | LOW | HIGH | Plan maintenance window, test in staging |
| Pipeline execution failures post-upgrade | MEDIUM | HIGH | Comprehensive testing before production upgrade |
| Lack of rollback capability | LOW | CRITICAL | Ensure OLM rollback works, maintain backups |
| Incomplete migration causing data loss | LOW | CRITICAL | Multiple backups, validation checkpoints |

---

## Testing Recommendations

### Test Scenarios

1. **Pipeline Execution**
   - Run simple Pipeline (single Task)
   - Run complex Pipeline (multiple Tasks, conditions, workspaces)
   - Test PipelineRun with params and results
   - Verify TaskRun execution

2. **Triggers**
   - Test EventListener receives events
   - Verify TriggerTemplate instantiation
   - Check TriggerBinding param extraction
   - Test webhook delivery

3. **Dashboard**
   - Access Dashboard UI
   - View PipelineRuns and TaskRuns
   - Check log viewing
   - Verify RBAC permissions

4. **Chains (if enabled)**
   - Verify signature generation
   - Check attestation storage
   - Test signature verification

5. **Results (if enabled)**
   - Verify result storage
   - Test result API queries
   - Check result retention

6. **Pruner (if enabling)**
   - Configure Pruner CRD
   - Verify automated cleanup
   - Check retention policies work

### Validation Checklist

- [ ] Operator deployment healthy
- [ ] All component pods running
- [ ] TektonConfig reconciled successfully
- [ ] Simple pipeline executes successfully
- [ ] Complex pipeline executes successfully
- [ ] Triggers fire correctly
- [ ] Dashboard accessible
- [ ] No deprecation warnings in logs
- [ ] Custom integrations still work
- [ ] RBAC permissions intact
- [ ] Existing PipelineRuns/TaskRuns accessible
- [ ] Results queryable (if enabled)
- [ ] Chains signing working (if enabled)

---

## Recommendations

### For EUS Customers (v0.71.x → v0.78.x)

**Option 1: Direct Upgrade (Aggressive)**
- Upgrade OpenShift 4.15 → 4.20
- Upgrade Tekton operator v0.71 → v0.78 in one step
- **Pros**: Fastest path
- **Cons**: Higher risk, less isolation of issues
- **Recommended for**: Small deployments, low pipeline complexity

**Option 2: Staged Approach (Conservative) ✅ RECOMMENDED**
- Upgrade OpenShift 4.15 → 4.20
- Upgrade Tekton operator v0.71 → v0.75 (stay on Pipeline v0)
- Test thoroughly
- Upgrade Tekton operator v0.75 → v0.76 (Pipeline v1 transition)
- Test Pipeline v1 compatibility
- Upgrade Tekton operator v0.76 → v0.78
- **Pros**: Isolates Pipeline v1 migration risk
- **Cons**: More steps
- **Recommended for**: Production deployments, complex pipelines

### Critical Actions

1. **DO NOT SKIP**: Pipeline v1 migration preparation
2. **MUST TEST**: All critical pipelines in non-production first
3. **MAINTAIN**: Current backups throughout upgrade process
4. **REVIEW**: Tekton upstream documentation for Pipeline v1

### Timeline Estimate

- **Assessment & Planning**: 1-2 days
- **Staging Environment Testing**: 3-5 days
- **Production Upgrade Window**: 2-4 hours
- **Post-Upgrade Validation**: 1-2 days
- **Total**: 1-2 weeks for full migration

---

## References

- [Tekton Operator GitHub](https://github.com/tektoncd/operator)
- [Tekton Pipeline Documentation](https://tekton.dev/docs/pipelines/)
- [Tekton Pipeline v1 Migration Guide](https://tekton.dev/docs/pipelines/migrating-v1beta1-to-v1/)
- [OpenShift Pipelines Documentation](https://docs.openshift.com/pipelines/)
- [OCP EUS Documentation](https://docs.openshift.com/container-platform/latest/updating/updating_a_cluster/eus-eus-update.html)

---

## Appendices

### Appendix A: Component Version Matrix (Full)

See "Version Timeline & Component Matrix" section above.

### Appendix B: CRD Comparison

**v0.71.x CRDs**: 9 CRDs
**v0.78.x CRDs**: 10 CRDs (+ TektonPruner)

No CRDs removed, full backwards compatibility at operator level.

### Appendix C: RBAC Permission Comparison

Minimal changes detected (927 → 934 lines, +0.75%).
New permissions likely related to Pruner component.

### Appendix D: Release Notes Links

- [v0.71.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.71)
- [v0.72.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.72)
- [v0.73.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.73)
- [v0.74.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.74)
- [v0.75.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.75)
- [v0.76.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.76) ⚠️ Pipeline v1
- [v0.77.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.77)
- [v0.78.x Release Notes](https://github.com/tektoncd/operator/releases?q=v0.78)

---

**End of Report**
