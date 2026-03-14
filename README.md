# teleport-knowledge-graph

Teleport Cluster Digital Twin staged build repository.

## Current status
- `MASTER_PROMPT.md`: full project specification prompt.
- `PHASE1_ARCHITECTURE.md`: architecture and contracts.
- **Phase 2 scaffolding implemented**: module, entrypoint, config/logging, interfaces, domain model, graph package, storage abstractions.
- **Phase 3 integration implemented**: Teleport auth/client factory contracts, concrete collectors, collector orchestration with partial-failure handling, and snapshot assembly.

## Project layout (implemented so far)
- `cmd/twin`: entrypoint bootstrap + collection execution + snapshot persistence.
- `internal/config`: configuration model and validation.
- `internal/logging`: structured logger creation.
- `internal/teleport`: auth loaders, client factory, error kind mapping, resource interfaces.
- `internal/collectors`: collector interfaces, concrete cluster/users/roles/nodes collectors, orchestrator.
- `internal/model`: typed entities, relationships, snapshot envelope, deterministic IDs.
- `internal/graph`: typed in-memory graph + traversal + path search + export.
- `internal/storage`: filesystem JSON snapshot store and in-memory current snapshot store.
- `internal/analysis`, `internal/diff`, `internal/api`, `internal/cli`, `internal/ai`: foundational contracts and early services.
- `pkg/twinapi`: schema version contract.

## Quick start
```bash
go build ./...
go test ./...
go run ./cmd/twin
```
