```yaml
version: "1.0"
description: "Core logic for processing customer transactions"

tables:
  - name: "customer_orders"
    description: "Tracks all successful checkout events"
    columns:
      - name: "id"
        type: "uuid"
        constraints: ["primary_key", "not_null"]
      - name: "order_value"
        type: "decimal"
        constraints: ["not_null"]
```

```sql
CREATE TABLE IF NOT EXISTS customer_orders (
  id  UUID PRIMARY KEY
  order_value DECIMAL NOT NULL
);
```

```yaml
tables:
  name: users
  columns:
    - name: id
      type: integer
      primary_key: true
      unique: true
      nullable: false
      default: autoincrement
      foreign_key: null
      example: 1
      description: Primary key, unique identifier for users table
    - name: username
      primary_key: false
      type: string
      unique: true
      default: false
      example: "jsmith"
      description: User's log in name
      nullable: false
    - name: email
      type: string
      unique: true
      foreign_key: false
      example: "jsmith@example.com"
      description: User's email address
    - name: create_at
      type: datetime
      default: CURRENT_TIMESTAMP
      primary_key: false
      unique: false
      nullable: false
      example: "YYYY-MM-DD HH24:MI:SS.FF or 2024-04-17 15:41:38.29"
      description: When the user was created
  actions:
    - create
    - read
    - update
    - delete
```
## Validation errors
1. Validation will have layers of complexity:
2. Structural — does the YAML have the right shape? (Parser handles this)
3. Identity — are names unique within their scope? (Duplicate detection)
4. Referential — do FKs and references resolve? (Cross-table lookup)
5. Type — are FK column types compatible? (Type compatibility)
6. Semantic — does the design make sense? (e.g., circular FK references, orphan tables)
7. Normalisation — does the schema follow 1NF/2NF/3NF? (Higher-level rules)
8. Style/lint — naming conventions, reserved words, etc. (Cosmetic but useful)

## New YAML Contract format
Probs a better representationof how constraints actually work and should mean more consistantly usable SQL code for generatation
```yaml
version: "1.0"
description: "Core schema for the orders platform"

tables:
  - name: customers
    description: "Customer master data"
    columns:
      - name: id
        type: uuid
        constraints: ["not_null"]
      - name: email
        type: varchar
        constraints: ["not_null"]
      - name: name
        type: varchar
        constraints: ["not_null"]
    primary_key: [id]
    unique:
      - [email]

  - name: orders
    description: "Customer orders"
    columns:
      - name: id
        type: uuid
        constraints: ["not_null"]
      - name: customer_id
        type: uuid
        constraints: ["not_null"]
        references:
          table: customers
          column: id
      - name: order_total
        type: decimal
        constraints: ["not_null"]
    primary_key: [id]

  - name: order_items
    description: "Line items for each order"
    columns:
      - name: order_id
        type: uuid
        constraints: ["not_null"]
        references:
          table: orders
          column: id
      - name: product_id
        type: uuid
        constraints: ["not_null"]
        references:
          table: products
          column: id
      - name: quantity
        type: integer
        constraints: ["not_null"]
    primary_key: [order_id, product_id]   # composite PK
```

## Next things to work on

- Profiling generator (port Python → Postgres SQL, reads JSON contract)
- Validation generator (derived rules from constraints, also reads JSON)
- Synthetic data MVP (reads JSON, uses topological order, generates rows with FK integrity)
- deployed on Fly
- Write-up: README / blog post explaining the polyglot pipeline

End state: working demo site showing YAML in → DDL + docs + diagram + validation SQL + synthetic data out, with a real architectural narrative behind it.

Next 3 reasonable choices:

- JSON enrichment. Connects the spec-and-consumers framing to actual code — JSON output starts being something a separate tool would actually want.

## Future architecture: per-generator validation requirements

The validator currently produces issues with a single global severity. As more
generators are added with different tolerance for the same issues (e.g., FK
cycles are fine for DDL but break synthetic data generation), the architecture
needs to evolve.

### Proposed model

Separate analysis from interpretation:

1. **Validator** produces a structured ValidationResult containing:
   - Which checks passed (map[CheckCode]bool)
   - The underlying issues (for messaging)
   
2. **Generators** declare their requirements (a set of CheckCodes that must pass)
   and consume the ValidationResult to decide:
   - Whether they can run (all required checks passed)
   - Which non-required failures should be surfaced as warnings

3. **UI** displays:
   - Global issues from the validator
   - Per-generator status (ran with warnings, blocked by required check)

### Why this design

Validator is purely analytical — no interpretation of impact.
Generators own their requirements — colocated with the code that has them.
UI handles presentation — no business logic mixed with rendering.

### When to implement

When the second generator with different requirements is added. Currently:
- DDL, Docs, Mermaid all tolerate everything that isn't a hard validation error
- Single severity per check is sufficient
- Hypothetical synthetic data / migration generators would have different needs

## v1.1 polish ideas

### Streaming generation via goroutines + SSE
- Spawn each generator in its own goroutine
- Stream results to client as they complete via Server-Sent Events
- Mostly demo value (generators are fast); demonstrates concurrency
- Blog post material

### Auto-generate with last-good preservation
- Debounced regeneration on textarea input
- Successful generators update; failed ones keep last-good output
- Issues panel updates in real-time
- Better UX than button-triggered, but more complex to implement

## v1.1: Schema introspection

Read existing Postgres database, produce a YAML contract matching its schema.

### Value
- Closes the loop: contract → SQL → contract
- Demonstrates database introspection skills (information_schema, pg_catalog)
- Enables schema drift detection (compare stored contract to live DB)
- Validates the architectural choice of Schema as central type — both YAML parsing and DB introspection converge on the same internal representation

### Scope
- CLI only initially (security concerns with web UI accepting connection strings)
- Postgres only initially (other databases require per-database adapters)
- Common types only initially (arrays, ranges, JSONB, custom types: future)
- Skip exotic constraints initially (check, deferred, exclusion: future)

### Technical shape
- New `internal/introspect/` package
- Queries information_schema and pg_catalog
- Returns *Schema (existing type)
- CLI: `cgen introspect --db postgres://... --output contract.yaml`

## v1.1 headline feature: Live validation

Add continuous validation as the user types, similar to IDE linting.

### Implementation
- New endpoint: `POST /validate` returning issues without generation
- HTMX on textarea: `hx-post="/validate" hx-trigger="keyup changed delay:500ms"`
- Generate button disabled when errors exist
- Generate endpoint still validates server-side (defence in depth)

### Why
- Tighter feedback loop for user iterating on YAML
- "Soft linting" — same pattern as IDE error highlighting
- Blocks generation when state is invalid (no wasted generator runs)

### Blog post angle
"Adding live validation to ContractGen: how the tiered validator pattern made
this easy" — leveraging the architectural work in v1 to ship a UX improvement
in v1.1.

cli commands
```bash
# Defualt behaviour - generate all 3
go run ./cmd/cgen --input contracts/example.yml

# Just docs
go run ./cmd/cgen --input contracts/example.yml --doc

# Just DDL
go run ./cmd/cgen --input contracts/example.yml --ddl

# Docs and JSON
go run ./cmd/cgen --input contracts/example.yml --doc --json


# Customt output dir
mkdir test-output
go run ./cmd/cgen --input contracts/example.yml --output-dir test-output

# missing input should fialt, "--input required"
go run ./cmd/cgen

# bad contract should fail, print all validation errors
go run ./cmd/cgen --input contracts/badExample.yml


# see help
go run ./cmd/cgen --help

# build a binary instead
go build -o cgen.exe ./cmd/cgen
./cgen.exe --input contracts/example.yml
```

```bash
# clean up
 rm schema.sql schema.md schema.json
 ```

 testing 
 ```bash
  go test -coverprofile=coverage.out ./internal/schema && go tool cover -html=coverage.out
  ```
  ```bash
go test -v ./internal/schema
  ```

 standard flow
 yaml -> parse -> validate ->  build schema -> (topo ->) generator 