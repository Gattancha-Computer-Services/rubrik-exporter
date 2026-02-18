# Rubrik Exporter GraphQL Migration Plan

## Overview
Rubrik CDM has migrated from REST API to GraphQL. This document outlines the migration strategy.

## Current State
- **16 REST API endpoints** currently used
- **Authentication**: Basic auth + session tokens
- **Data format**: JSON responses with consistent structure

## Migration Strategy

### Phase 1: Infrastructure Setup
1. ✅ **Backup created** - All current files backed up
2. **Add GraphQL dependencies**:
   ```bash
   go get github.com/machinebox/graphql
   ```
3. **Create GraphQL client** - New `GraphQLClient` struct
4. **Update authentication** - GraphQL uses Bearer tokens

### Phase 2: API Migration (By Priority)

#### High Priority (Core Metrics)
1. **System Storage** (`/api/internal/stats/system_storage`) ✅ **MIGRATED**
   - GraphQL: `system { storage { total, used, available, snapshot, liveMount, miscellaneous } }`
   - Status: Ready for implementation

2. **Nodes** (`/api/internal/node`) ✅ **MIGRATED**
   - GraphQL: `nodes { id, name, status, ipAddress, cluster { id, name } }`
   - Status: Ready for implementation

3. **VMware VMs** (`/api/v1/vmware/vm`) ✅ **MIGRATED**
   - GraphQL: `vmwareVms { edges { node { id, name, effectiveSlaDomain { id, name } } } }`
   - Status: Ready for implementation

#### Medium Priority (Additional Metrics)
4. **Archive Locations** (`/api/internal/archive/location`) ✅ **MIGRATED**
5. **Managed Volumes** (`/api/internal/managed_volume`) ✅ **MIGRATED**
6. **Statistics** (multiple `/api/internal/stats/*` endpoints) ✅ **MIGRATED**
   - Per VM Storage ✅
   - Streams Count ✅
   - Data Location Usage ✅
   - Physical Ingest Time Series ✅
   - Archival Bandwidth Time Series ✅
   - Runway Remaining ✅
   - Average Storage Growth ✅
7. **Reports** (`/api/internal/report`) ✅ **MIGRATED**

#### Low Priority (Nice-to-have)
8. **Nutanix VMs** (`/api/internal/nutanix/vm`) ✅ **MIGRATED**
9. **Hyper-V VMs** (`/api/internal/hyperv/vm`) ✅ **MIGRATED**

### Phase 3: Testing & Validation
1. **Unit tests** for each GraphQL query
2. **Integration tests** against Rubrik CDM
3. **Metrics validation** - Ensure all Prometheus metrics still work
4. **Performance testing** - GraphQL should be more efficient

## Migration Progress

### ✅ Completed
1. **Infrastructure Setup**
   - ✅ Added GraphQL dependencies to go.mod
   - ✅ Created GraphQL client (graphql_client.go)
   - ✅ Updated Rubrik struct to include GraphQL client
   - ✅ Updated authentication to initialize GraphQL client

2. **API Migrations (ALL 16 ENDPOINTS COMPLETED)**
   - ✅ **System Storage** (`/api/internal/stats/system_storage`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `GetSystemStorage()`
   - ✅ **Nodes** (`/api/internal/node`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `GetNodes()`
   - ✅ **VMware VMs** (`/api/v1/vmware/vm`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `ListVmwareVM()`
   - ✅ **Nutanix VMs** (`/api/internal/nutanix/vm`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `ListNutanixVM()`
   - ✅ **Hyper-V VMs** (`/api/internal/hyperv/vm`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `ListHypervVM()`
   - ✅ **Archive Locations** (`/api/internal/archive/location`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `GetArchiveLocations()`
   - ✅ **Managed Volumes** (`/api/internal/managed_volume`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `GetManagedVolumes()`
   - ✅ **Statistics APIs** (7 endpoints)
     - Per VM Storage: `GetPerVMStorage()`
     - Streams Count: `GetStreamCount()`
     - Data Location Usage: `GetDataLocationUsage()`
     - Physical Ingest Time Series: `GetPhysicalIngest()`
     - Archival Bandwidth Time Series: `GetArchivalBandwith()`
     - Runway Remaining: `GetRunawayRemaining()`
     - Average Storage Growth: `GetAverageStorageGrowthPerDay()`
   - ✅ **Reports** (`/api/internal/report`)
     - GraphQL query implemented
     - Fallback to REST API added
     - Method: `GetReports()`

### 🔄 Ready for Testing
3. **Phase 3: Testing & Validation**
   - ⏳ Unit tests for each GraphQL query
   - ⏳ Integration tests against Rubrik CDM
   - ⏳ Metrics validation - Ensure all Prometheus metrics still work
   - ⏳ Performance testing - GraphQL should be more efficient

## Benefits of GraphQL Migration

1. **Single endpoint** - All queries go to `/api/graphql`
2. **Precise data fetching** - Request only needed fields
3. **Better performance** - Reduced over/under-fetching
4. **Type safety** - GraphQL schemas provide better validation
5. **Future-proof** - GraphQL is Rubrik's API direction

## Risks & Mitigations

### Risk: API Changes
- **Mitigation**: Test against multiple Rubrik versions
- **Fallback**: Keep REST client as backup option

### Risk: Performance Impact
- **Mitigation**: GraphQL should be faster (single connection, precise queries)
- **Monitoring**: Add query timing metrics

### Risk: Breaking Changes
- **Mitigation**: Comprehensive testing, gradual rollout
- **Rollback**: Backup allows quick reversion

## Timeline Estimate

- **Phase 1 (Infrastructure)**: ✅ 1-2 days (COMPLETED)
- **Phase 2 (Core APIs)**: ✅ 3-5 days (COMPLETED - All 16 APIs migrated)
- **Phase 3 (Testing)**: 📋 2-3 days (READY FOR TESTING)
- **Total**: 1-2 weeks (Expected completion: 1 week remaining)

## Testing Strategy

1. **Unit Tests**: Mock GraphQL responses
2. **Integration Tests**: Real Rubrik CDM instance
3. **Load Tests**: Multiple concurrent queries
4. **Regression Tests**: All existing metrics still work

## Rollback Plan

If issues arise:
1. Restore from backup: `Copy-Item -Recurse backup-*\* . -Force`
2. Revert to REST API implementation
3. Document issues for future GraphQL migration

## Resources

- [Rubrik GraphQL API Documentation](https://rubrikinc.github.io/rubrik-api-docs/)
- [GraphQL Go Client](https://github.com/machinebox/graphql)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)

## Next Steps

1. **TEST IMMEDIATELY**: Deploy and test against a real Rubrik CDM instance
2. **Validate All Metrics**: Ensure all Prometheus metrics work with GraphQL
3. **Performance Comparison**: Compare GraphQL vs REST API performance
4. **Gradual Rollout**: Start with GraphQL, keep REST as fallback during transition
5. **Remove REST Code**: Once GraphQL is stable, remove REST API fallbacks
6. **Update Documentation**: Update README with GraphQL-specific requirements