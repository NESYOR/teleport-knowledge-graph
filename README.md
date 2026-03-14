# teleport-knowledge-graph

Teleport Cluster Digital Twin staged build repository.

## Current status
- `MASTER_PROMPT.md`: full project specification prompt.
- `PHASE1_ARCHITECTURE.md`: architecture and contracts.
- **Phase 2 scaffolding implemented**: module, entrypoint, config/logging, interfaces, domain model, graph package, storage abstractions.

## Project layout (implemented in this phase)
- `cmd/twin`: entrypoint bootstrap.
- `internal/config`: configuration model and validation.
- `internal/logging`: structured logger creation.
- `internal/model`: typed entities, relationships, snapshot envelope, deterministic IDs.
- `internal/graph`: typed in-memory graph + traversal + path search + export.
- `internal/storage`: filesystem JSON snapshot store and in-memory current snapshot store.
- `internal/teleport`, `internal/collectors`, `internal/analysis`, `internal/diff`, `internal/api`, `internal/cli`, `internal/ai`: foundational interfaces/contracts.
- `pkg/twinapi`: schema version contract.

## Quick start
```bash
go build ./...
```
