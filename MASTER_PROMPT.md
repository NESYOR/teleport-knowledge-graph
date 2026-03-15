# Teleport Cluster Digital Twin — Codex Master Prompt

Use the following master prompt when you want Codex to generate the full Teleport Cluster Digital Twin project specification.

## Master prompt

```text
You are a senior Go systems engineer, security architect, and AI infrastructure engineer.

Build a complete production-quality project called:

Teleport Cluster Digital Twin

Goal:
Create a Go-based platform that connects to a Teleport cluster using the official Go Teleport API client and builds an AI-ready “digital twin” of the cluster. The system must collect cluster state, normalize it into a graph-based infrastructure model, expose queryable APIs, and support AI reasoning over access relationships, resource topology, and security posture.

High-level concept:
Teleport Cluster -> Go collector -> normalized graph model -> analysis engine -> AI query layer

Core problem to solve:
Existing Teleport tools expose raw resources and audit data, but they do not provide a unified digital twin that models:
- identities
- roles
- traits
- users
- nodes
- databases
- kubernetes clusters
- applications
- windows desktops
- trusted clusters
- access paths
- privilege inheritance
- resource visibility relationships
- security risks

This project must create that digital twin and make it usable by:
1. humans via CLI and REST API
2. AI systems via structured JSON schemas and tool-friendly outputs

Primary outcomes:
1. Discover all accessible cluster information through the Teleport Go API client
2. Build a normalized in-memory graph of the cluster
3. Persist snapshots as JSON
4. Provide a query layer for exploring relationships
5. Provide an analysis layer for privilege and exposure insights
6. Provide an AI-ready interface for natural language and structured reasoning integrations
7. Be modular, extensible, testable, and production-ready

Technical requirements:
- Language: Go
- Go version: latest stable supported release
- Package management: Go modules
- Follow idiomatic Go project structure
- Use the Teleport Go API client from:
  github.com/gravitational/teleport/api/client
- Build a real working project, not pseudo-code
- Include complete code, README, configs, tests, Makefile, and examples
- No stubs unless clearly marked and justified
- All exported types/functions must have clear comments
- Use structured logging
- Use context cancellation everywhere appropriate
- Handle partial failures gracefully
- Prefer interfaces around Teleport client interactions for testability
- Output must be compile-ready

Project vision:
This is not just an inventory exporter.
It is an infrastructure knowledge system and AI-ready security reasoning engine for Teleport.

What the system must do

A. Cluster collection
Collect as much cluster information as the authenticated Teleport identity can access, including at minimum:
- cluster name
- current user
- roles visible to the caller
- users visible to the caller if available
- nodes
- database servers / databases
- kubernetes servers / kube clusters
- applications if accessible
- windows desktops
- trusted clusters
- cluster networking config if available
- auth preference / session recording / relevant security settings if accessible
- server info or version metadata where useful
- labels, traits, static and dynamic metadata when available

Important:
Design collectors so that lack of permission for one resource type does not fail the full snapshot. The system should return partial success with errors grouped by collector.

B. Digital twin graph model
Convert all collected resources into a unified graph representation.

Model the following entity types:
- Cluster
- User
- Role
- Trait
- Node
- Database
- DatabaseServer
- KubernetesCluster
- KubernetesServer
- Application
- WindowsDesktop
- TrustedCluster
- Label
- Login
- Namespace
- Principal
- ResourceSelector
- AccessCapability
- Rule
- Condition
- Edge / Relationship

Model the following relationship types:
- USER_HAS_ROLE
- ROLE_ALLOWS_LOGIN
- ROLE_MATCHES_LABEL
- USER_CAN_ACCESS_NODE
- USER_CAN_ACCESS_DB
- USER_CAN_ACCESS_KUBE
- USER_CAN_ACCESS_DESKTOP
- ROLE_GRANTS_VERB
- RESOURCE_IN_CLUSTER
- TRUSTS_CLUSTER
- RESOURCE_HAS_LABEL
- USER_HAS_TRAIT
- ROLE_REQUIRES_TRAIT
- ACCESS_PATH
- POSSIBLE_PRIVILEGE_PATH

The graph must support both:
1. in-memory operation
2. serialization to JSON for persistence/export

C. Snapshot engine
Implement a snapshot engine that:
- performs a full collection run
- records collection metadata
- stores timestamps
- records per-collector errors
- computes deterministic IDs for nodes/entities where appropriate
- writes normalized snapshots to disk as JSON
- supports loading snapshots back into memory

D. Analysis engine
Implement analysis over the graph.

At minimum support:
1. User access analysis
   - Which resources can a user access?
   - Through which roles?
   - Why?

2. Resource exposure analysis
   - Which users or roles can access a resource?
   - Which roles expose prod resources?
   - Which resources are broadly exposed?

3. Risk analysis
   - wildcard labels
   - broad logins like root
   - dangerous access combinations
   - roles with excessive reach
   - cross-cluster trust risk indicators
   - possible privilege expansion paths

4. Topology analysis
   - resource counts by type
   - cluster composition
   - trust relationships
   - namespace or environment segmentation

5. Diff analysis
   Support comparing two snapshots and reporting:
   - added resources
   - removed resources
   - modified roles
   - changed label scopes
   - changed access reach

E. Query interfaces
Provide 3 interfaces:

1. CLI
Commands like:
- twin collect
- twin show summary
- twin show users
- twin show roles
- twin show resources
- twin explain user <username>
- twin explain resource <resource-id>
- twin analyze risks
- twin diff old.json new.json

2. REST API
Provide endpoints such as:
- GET /healthz
- POST /collect
- GET /snapshot/current
- GET /summary
- GET /users/:name/access
- GET /resources/:id/exposure
- GET /roles/:name
- GET /analysis/risks
- POST /diff

3. AI tool-friendly API
Expose structured JSON responses specifically designed for LLM and agent use.
Responses must be deterministic, concise, and schema-stable.

F. AI integration layer
Do not hardcode to one LLM provider.
Design an abstraction layer for AI reasoning.

Implement:
- interfaces for prompt generation
- structured context builders from graph subgraphs
- AI-ready summarizers for users, roles, and resources
- functions that produce compact graph context for a question

Examples of supported future AI questions:
- “Why can user alice access production-db?”
- “Which roles expose root access to production nodes?”
- “What changed in access posture between snapshot A and B?”
- “What is the shortest privilege path from user X to resource Y?”
- “Which clusters trust this cluster and what resources are exposed through that?”

Include a local rule-based reasoning layer first.
Then provide pluggable AI adapter interfaces for future LLM providers.

G. Extensibility
Design the codebase so future features can be added easily:
- graph database backends like Neo4j
- MCP server support
- LLM providers
- audit-log ingestion
- periodic collection
- web UI
- policy recommendation engine

Architecture requirements

Use a clean layered architecture.

Suggested top-level structure:
- cmd/
- internal/config/
- internal/logging/
- internal/teleport/
- internal/collectors/
- internal/model/
- internal/graph/
- internal/analysis/
- internal/diff/
- internal/api/
- internal/cli/
- internal/storage/
- internal/ai/
- pkg/
- examples/
- test/

Mandatory design expectations:
1. Separate Teleport resource collection from graph modeling
2. Separate graph modeling from analysis
3. Separate analysis from transport layers
4. Use interfaces for client, storage, and AI adapters
5. Make collectors independently testable
6. Use dependency injection where appropriate
7. Avoid global state
8. Prefer composition over inheritance-like patterns

Implementation details

Teleport integration:
- Build a Teleport client wrapper
- Support auth via tsh profile if available
- Support auth via identity file
- Support config-driven connection options
- Encapsulate Teleport client initialization behind a factory
- Gracefully handle unavailable methods or permission denials

Graph design:
Implement a first-class graph package that supports:
- nodes/entities
- typed edges
- lookup by ID
- lookup by type
- adjacency traversal
- reverse traversal
- subgraph extraction
- shortest path or BFS-based path discovery for access explanations

Important:
The graph should be domain-aware, not a generic untyped map soup.

Data modeling:
Define strong typed domain models.
Avoid excessive use of map[string]interface{}.
Use typed structs with explicit JSON tags.

Error handling:
- Return partial collection results
- Record detailed collector errors
- Differentiate permission errors, network errors, parse errors, and unsupported-method errors
- Never crash on one failing collector unless it prevents basic operation

Testing:
Include:
- unit tests
- collector tests using mocks
- graph tests
- analysis tests
- diff tests
- API tests
- CLI smoke tests if practical

Testing should cover:
- empty cluster cases
- partial permission cases
- broad role access
- trust relationships
- diff of changed role selectors
- user access explanation logic

Developer experience:
Include:
- README with architecture diagram in markdown text
- example commands
- example config file
- sample snapshot JSON
- Makefile targets
- lint/test/build instructions

Non-functional requirements:
- readable code
- clear naming
- no unnecessary abstractions
- no placeholder enterprise fluff
- make it feel like a real serious infra/security project
- optimize for maintainability and extension

Output expectations:
Generate the complete project in phases.

Phase 1:
- project plan
- architecture explanation
- final directory structure
- domain model description
- graph model description
- API contract outline
- CLI command design

Phase 2:
- full source code for all files
- include go.mod and go.sum where possible
- README
- example config
- Makefile

Phase 3:
- tests
- sample snapshot data
- usage walkthrough
- future roadmap

Important generation rules:
- Do not skip files
- Do not collapse major files with “rest omitted”
- Do not provide pseudo-code where real code is expected
- When code is lengthy, still provide complete code
- Ensure imports are coherent
- Ensure the project compiles logically
- Favor correctness over brevity

Additional advanced features to include if feasible:
- configurable collection timeout
- per-collector enable/disable flags
- snapshot metadata including auth source and cluster identity
- graph export in JSON
- explain endpoint that returns path-based reasoning
- risk scoring model for roles/resources
- access summary grouped by environment labels
- snapshot diff severity classification

Important product framing:
This project should look like something a senior platform/security engineer or AI infra engineer would build for:
- Teleport security posture reasoning
- AI copilot integrations
- Zero Trust access visualization
- production access analysis
- change intelligence and drift detection

Code style preferences:
- idiomatic Go
- small focused packages
- table-driven tests where appropriate
- use stdlib where possible
- minimize unnecessary external dependencies
- use chi or net/http for REST API, choose one and be consistent
- use cobra for CLI if helpful, otherwise build a clean stdlib CLI

Deliverable format:
First output the architecture and project structure.
Then output all files one by one in clear order with file paths.
Ensure each file is complete.

Also include:
1. a list of assumptions made about Teleport API availability
2. explicit extension points for future LLM or graph DB integration
3. a short section called “How this becomes an AI reasoning system”

Now begin with Phase 1 only.
```

## Recommended staged prompting sequence

### Prompt 1 — architecture

```text
Using the master spec I gave you, complete only Phase 1:
- architecture
- project structure
- domain model
- graph model
- API contract
- CLI design
Do not generate code yet.
Be concrete and opinionated.
```

### Prompt 2 — scaffolding

```text
Using the approved architecture, generate the full repository scaffolding and all foundational files:
- go.mod
- main entrypoint
- config
- logging
- interfaces
- domain models
- graph package
- storage abstractions
- README
No TODO placeholders unless unavoidable.
```

### Prompt 3 — Teleport integration

```text
Generate the Teleport integration layer:
- client factory
- auth/profile loading
- collectors
- error handling
- collector orchestration
- snapshot assembly
Use interfaces for testability.
```

### Prompt 4 — analysis engine

```text
Generate the graph analysis and reasoning engine:
- user access explanation
- resource exposure analysis
- risk analysis
- path traversal
- snapshot diff engine
Include tests.
```

### Prompt 5 — interfaces

```text
Generate the CLI and REST API layers on top of the existing packages.
Implement all routes and commands from the architecture.
Include request/response models and tests.
```

### Prompt 6 — AI layer

```text
Generate the AI-ready context builder and reasoning adapter interfaces:
- compact graph context generation
- deterministic summarizers
- LLM adapter abstraction
- prompt templates
- structured output contracts
Do not bind to a single LLM provider.
```

### Prompt 7 — hardening

```text
Review the full project and improve it for:
- compile correctness
- consistency
- missing imports
- package cohesion
- error handling
- test quality
- README clarity
Then output a final repository tree.
```

Best practice: generate **Phase 1 first**, then iterate. If you ask for the whole project in one shot, it may drift, duplicate types, or create mismatched package boundaries.
