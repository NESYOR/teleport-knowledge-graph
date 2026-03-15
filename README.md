# teleport-knowledge-graph

Teleport Cluster Digital Twin staged build repository.

## Current status
- `MASTER_PROMPT.md`: full project specification prompt.
- `PHASE1_ARCHITECTURE.md`: architecture and contracts.
- **Phase 2 scaffolding implemented**: module, entrypoint, config/logging, interfaces, domain model, graph package, storage abstractions.
- **Phase 3 integration implemented**: Teleport auth/client factory contracts, concrete collectors, collector orchestration with partial-failure handling, and snapshot assembly.
- **Phase 4 analysis implemented**: user access explanation, resource exposure analysis, deterministic risk checks, topology summary, and enhanced snapshot diff classification.
- **Phase 5 interfaces implemented**: CLI command runner and REST API server with collection, summary, explain, risk, role, and diff endpoints.

## Quick start
```bash
go build ./...
go test ./...

# CLI
./teleport-cluster-digital-twin collect
./teleport-cluster-digital-twin show summary

# REST API
./teleport-cluster-digital-twin serve
curl -s localhost:8080/healthz
```
