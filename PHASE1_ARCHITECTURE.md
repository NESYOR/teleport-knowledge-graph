# Teleport Cluster Digital Twin — Phase 1 Architecture

This document provides **Phase 1 only**:
- project plan
- architecture explanation
- final directory structure
- domain model
- graph model
- API contract outline
- CLI command design

No implementation code is included in this phase.

---

## 1) Project plan (opinionated)

### Delivery sequence
1. **Foundation and contracts**
   - Freeze package boundaries, interfaces, IDs, and JSON schemas.
   - Define snapshot envelope and error taxonomy.
2. **Telemetry ingestion and normalization**
   - Add Teleport client factory + collectors with per-collector fault isolation.
   - Convert raw resources to typed domain entities.
3. **Graph build + query primitives**
   - Build typed graph with forward/reverse edges, indexed lookup, BFS pathing.
4. **Analysis and diff engine**
   - User access explanation, resource exposure, risks, topology, and snapshot diff.
5. **Transport layers**
   - CLI and REST APIs using the same internal services.
6. **AI-ready layer**
   - Deterministic summarizers and compact context builders.
7. **Hardening**
   - Test expansion, schema checks, docs, and operational polish.

### Milestones
- **M1:** Snapshot collection succeeds with partial failures and deterministic export.
- **M2:** Graph explains `user -> role -> selector -> resource` paths.
- **M3:** Risk and diff outputs are stable JSON contracts.
- **M4:** CLI + REST parity for all core operations.

---

## 2) Architecture

### Layered architecture

```text
Teleport Cluster
   │
   ▼
internal/teleport (client factory, auth, wrappers)
   │
   ▼
internal/collectors (resource-specific collectors, partial-failure orchestration)
   │
   ▼
internal/model (typed canonical entities)
   │
   ▼
internal/graph (typed nodes/edges, traversal, indexing)
   │
   ├──► internal/analysis (access/exposure/risk/topology)
   ├──► internal/diff (snapshot-to-snapshot change intelligence)
   └──► internal/ai (context builders, deterministic summarizers, adapters)

internal/storage (snapshot persistence)
internal/api (REST handlers)
internal/cli (command handlers)
cmd/twin (process entrypoint and wiring)
```

### Dependency rules
- `internal/teleport` depends only on config/logging and Teleport SDK.
- `internal/collectors` depends on teleport interfaces and model contracts.
- `internal/graph` depends on model only.
- `analysis`, `diff`, and `ai` depend on graph + model (not transport).
- `api` and `cli` depend on service interfaces only.
- `cmd/twin` composes dependencies via dependency injection.

### Runtime modes
- **One-shot CLI mode**: collect/analyze/export and exit.
- **Service mode**: API server keeps current snapshot in memory and supports collect/refresh.

---

## 3) Final directory structure

```text
teleport-cluster-digital-twin/
├── cmd/
│   └── twin/
│       └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── validate.go
│   ├── logging/
│   │   └── logger.go
│   ├── teleport/
│   │   ├── client_factory.go
│   │   ├── auth_profile.go
│   │   ├── auth_identity.go
│   │   ├── interfaces.go
│   │   └── errors.go
│   ├── collectors/
│   │   ├── orchestrator.go
│   │   ├── collector.go
│   │   ├── cluster_collector.go
│   │   ├── users_collector.go
│   │   ├── roles_collector.go
│   │   ├── nodes_collector.go
│   │   ├── db_collector.go
│   │   ├── kube_collector.go
│   │   ├── apps_collector.go
│   │   ├── desktops_collector.go
│   │   ├── trusts_collector.go
│   │   └── settings_collector.go
│   ├── model/
│   │   ├── types.go
│   │   ├── entities.go
│   │   ├── relationships.go
│   │   ├── snapshot.go
│   │   └── ids.go
│   ├── graph/
│   │   ├── graph.go
│   │   ├── index.go
│   │   ├── traversal.go
│   │   ├── path.go
│   │   └── export.go
│   ├── analysis/
│   │   ├── service.go
│   │   ├── access.go
│   │   ├── exposure.go
│   │   ├── risks.go
│   │   └── topology.go
│   ├── diff/
│   │   ├── service.go
│   │   └── severity.go
│   ├── storage/
│   │   ├── interface.go
│   │   ├── fs_json_store.go
│   │   └── current_store.go
│   ├── ai/
│   │   ├── contracts.go
│   │   ├── context_builder.go
│   │   ├── summarizer.go
│   │   └── rule_reasoner.go
│   ├── api/
│   │   ├── server.go
│   │   ├── routes.go
│   │   ├── handlers_collect.go
│   │   ├── handlers_summary.go
│   │   ├── handlers_explain.go
│   │   └── models.go
│   └── cli/
│       ├── app.go
│       ├── collect.go
│       ├── show.go
│       ├── explain.go
│       ├── analyze.go
│       └── diff.go
├── pkg/
│   └── twinapi/
│       └── schemas.go
├── examples/
│   ├── config.example.yaml
│   └── snapshot.sample.json
├── test/
│   ├── fixtures/
│   └── integration/
├── Makefile
├── go.mod
├── README.md
└── PHASE1_ARCHITECTURE.md
```

---

## 4) Domain model

### Entity types
- **Cluster**: identity for the source cluster, version, auth metadata.
- **User**: Teleport user, traits, role refs.
- **Role**: allow/deny sections, selectors, logins, rules.
- **Trait**: `key -> values` normalized as nodes for relationship reasoning.
- **Node**: SSH nodes with labels/namespace.
- **DatabaseServer** and **Database**: server process vs logical database target.
- **KubernetesServer** and **KubernetesCluster**: agent endpoint vs logical cluster.
- **Application**: app service targets.
- **WindowsDesktop**: desktop service resources.
- **TrustedCluster**: trust topology node.
- **Label**: normalized label key/value attached to resources.
- **Login / Principal / Namespace / ResourceSelector / AccessCapability / Rule / Condition** as first-class, queryable records.

### Snapshot envelope
- `snapshot_id` (deterministic hash from collection metadata + timestamp bucket)
- `collected_at`
- `source` (cluster name, auth source, caller)
- `collectors` (per-collector duration, status, counts)
- `collector_errors[]` (typed with kind and message)
- `entities[]`
- `relationships[]`
- `graph_metadata`

### Error taxonomy
- `permission_denied`
- `not_found`
- `unsupported_method`
- `network_error`
- `timeout`
- `parse_error`
- `internal_error`

Each collector must return either records, errors, or both.

---

## 5) Graph model

### Node contract
Each node must include:
- `id` (deterministic and type-scoped, e.g. `role:<cluster>:<name>`)
- `type`
- `name`
- `cluster`
- `attributes` (typed struct in memory, canonical JSON projection on export)

### Edge contract
- `id` (`<type>:<from>:<to>:<qualifier-hash>`)
- `type` (enum)
- `from`
- `to`
- `evidence` (rule/selector/trait source)
- `confidence` (rule-based certainty; default `1.0` for direct edges)

### Relationship types
- `USER_HAS_ROLE`
- `ROLE_ALLOWS_LOGIN`
- `ROLE_MATCHES_LABEL`
- `USER_CAN_ACCESS_NODE`
- `USER_CAN_ACCESS_DB`
- `USER_CAN_ACCESS_KUBE`
- `USER_CAN_ACCESS_DESKTOP`
- `ROLE_GRANTS_VERB`
- `RESOURCE_IN_CLUSTER`
- `TRUSTS_CLUSTER`
- `RESOURCE_HAS_LABEL`
- `USER_HAS_TRAIT`
- `ROLE_REQUIRES_TRAIT`
- `ACCESS_PATH`
- `POSSIBLE_PRIVILEGE_PATH`

### Required graph operations
- lookup by node ID
- lookup by node type
- outgoing/incoming adjacency
- constrained traversal by edge type
- subgraph extraction by seed + depth
- shortest path (BFS) with edge filters and explanation breadcrumbs

---

## 6) API contract outline (REST)

All responses are JSON and include:
- `request_id`
- `version`
- `data`
- optional `errors[]`

### Endpoints
1. `GET /healthz`
   - Returns service health and snapshot freshness.

2. `POST /collect`
   - Triggers collection.
   - Body options: timeout override, enabled collectors, persist flag.
   - Returns snapshot metadata + partial error summary.

3. `GET /snapshot/current`
   - Returns currently loaded snapshot envelope.

4. `GET /summary`
   - Returns counts by entity type, trust summary, exposure highlights.

5. `GET /users/{name}/access`
   - Returns accessible resources grouped by type.
   - Includes role and selector evidence.

6. `GET /resources/{id}/exposure`
   - Returns users/roles that can access the resource and why.

7. `GET /roles/{name}`
   - Returns normalized role model + computed reach stats.

8. `GET /analysis/risks`
   - Returns scored findings with severity and evidence paths.

9. `POST /diff`
   - Body: old snapshot ref + new snapshot ref (path or embedded).
   - Returns added/removed/modified entities and access posture changes with severity.

### Stability rules for AI clients
- deterministic ordering (sort by IDs and severity)
- stable enums for type/severity/status fields
- bounded output options (`limit`, `depth`, `include_evidence`)

---

## 7) CLI command design

Binary name: `twin`

### Commands
- `twin collect`
  - Collect from Teleport and optionally persist snapshot.
- `twin show summary`
  - Print topology/access summary from current or specified snapshot.
- `twin show users`
  - List users with role counts and reach counts.
- `twin show roles`
  - List roles with selector breadth and risk indicators.
- `twin show resources`
  - List resources by type/labels/namespace filters.
- `twin explain user <username>`
  - Explain access paths and evidence for a user.
- `twin explain resource <resource-id>`
  - Explain exposure paths for a resource.
- `twin analyze risks`
  - Emit risk findings; support `--format table|json`.
- `twin diff <old.json> <new.json>`
  - Compare snapshots and print severity-classified changes.

### Global flags
- `--config`
- `--output json|table`
- `--snapshot <path>`
- `--timeout`
- `--collector <name>` (repeatable)
- `--log-level`

---

## 8) Assumptions about Teleport API availability

1. Authenticated caller has access to at least cluster identity and one resource type.
2. Some list/read methods may fail with permission denied; collectors must continue.
3. Resource APIs differ by Teleport version; unsupported endpoints are mapped to `unsupported_method`.
4. `tsh` profile-based auth and identity-file auth can both be available, but either may be absent.
5. Labels and traits may be partially populated depending on backend and role permissions.

---

## 9) Explicit extension points

### LLM integration
- `AIAdapter` interface (provider-neutral): `Generate`, `Classify`, `StructuredAnswer`.
- Prompt/template registry isolated under `internal/ai`.
- Context builder contract that emits compact, bounded graph subgraphs.

### Graph DB integration
- `GraphStore` interface to persist/load graph from in-memory model.
- Initial implementation: in-memory + JSON snapshot.
- Future implementations: Neo4j, JanusGraph, or OpenSearch graph projections.

### Additional platform features
- periodic scheduler service (collection cron)
- audit-log ingestion pipeline
- MCP tool server wrapping analysis APIs
- policy recommendation module over risk findings

---

## 10) How this becomes an AI reasoning system

1. **Normalize first**: Teleport resources become canonical entities and typed edges.
2. **Reason locally**: deterministic rule-based analysis produces explainable paths.
3. **Context distillation**: AI layer extracts compact, relevant subgraphs.
4. **Provider abstraction**: external LLMs consume stable context and return structured outputs.
5. **Human + AI parity**: CLI/API and AI endpoints use the same analysis core, avoiding divergence.

This design ensures the system remains trustworthy for security analysis while becoming progressively more capable for AI-assisted investigation.
