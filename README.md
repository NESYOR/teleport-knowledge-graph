# Teleport Cluster Digital Twin

A Go-based platform that collects Teleport cluster state, normalizes it into a typed graph, and provides deterministic query/analysis surfaces for CLI, REST, and AI workflows.

## Architecture (text diagram)

```text
Teleport Cluster
   |
   v
internal/teleport (auth + client factory)
   |
   v
internal/collectors (resource collectors + orchestration)
   |
   v
internal/model (typed entities/relationships/snapshots)
   |
   v
internal/graph (typed graph + traversal/pathing)
   |\
   | +--> internal/analysis (access/exposure/risk/topology)
   | +--> internal/diff (snapshot-to-snapshot changes)
   | +--> internal/ai (context, summarizers, reasoner, adapters)
   |
   +--> internal/api (REST) / internal/cli (CLI)

internal/storage persists snapshots as JSON.
```

## Implemented phases

- **Phase 1:** Architecture and contracts (`MASTER_PROMPT.md`, `PHASE1_ARCHITECTURE.md`).
- **Phase 2:** Scaffolding, model, graph, storage, config, logging.
- **Phase 3:** Teleport integration contracts + collectors + orchestrator.
- **Phase 4:** Analysis engine and diff classification.
- **Phase 5:** CLI and REST interface layer.
- **Phase 6:** AI context builder, deterministic summarizers, rule reasoner, provider-neutral adapter contracts.
- **Phase 7:** Hardening for config parsing, deterministic exports, and runtime consistency.

## Quick start

```bash
make build
make test
```

### CLI flow

```bash
# collect a snapshot
./twin collect

# inspect summary from current snapshot
./twin show summary

# explain user access
./twin explain user alice

# analyze risks
./twin analyze risks
```

### REST flow

```bash
# start API
./twin serve

# health
curl -s localhost:8080/healthz

# collect
curl -s -X POST localhost:8080/collect

# summary
curl -s localhost:8080/summary

# risks
curl -s localhost:8080/analysis/risks
```

## Example artifacts

- Example config: `examples/config.example.yaml`
- Example snapshot: `examples/snapshot.sample.json`

## Developer targets

```bash
make fmt
make test
make build
make run
make serve
```

## Roadmap

1. Replace `LocalStaticClient` with real Teleport API client wiring.
2. Expand collectors to DB/Kubernetes/apps/desktops/trusted clusters/settings.
3. Add richer role selector evaluation and privilege-path analysis.
4. Add periodic collection scheduler and historical trend endpoints.
5. Add MCP server endpoints and pluggable graph database backend.
