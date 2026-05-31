# ContractGen

Generate SQL DDL, JSON, and documentation from YAML data contracts.

**[Try it live → contractgen.fly.dev](https://contractgen.fly.dev)**

Define your schema once in YAML; ContractGen parses, validates, and produces DDL, a data dictionary, an ER diagram, and JSON — keeping all of them in sync with a single source of truth.

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
    primary_key: [order_id, product_id]
```

## Try It

The fastest way to see what ContractGen does is the live web UI:

**[contractgen.fly.dev](https://contractgen.fly.dev)**

Paste a YAML contract and get DDL, docs, a Mermaid ER diagram, and validation
feedback rendered live. No install required.

To run it locally or use the CLI, see below.

## Installation

Requires Go 1.26+.

```bash
git clone https://github.com/Obiuno/contractgen
cd contractgen
go build -o contractgen ./cmd/cgen
```

Or run the CLI directly without building:

```bash
go run ./cmd/cgen --input contracts/example.yml
```

To run the web UI locally:

```bash
go run ./cmd/cgen-web
# then visit http://localhost:8080
```

## CLI Usage

```txt
--input        Path to YAML contract file (required)
--output-dir   Directory for output files (default: current directory)

Output flags (if none specified, all are produced):
--json     produce schema in JSON format
--ddl      produce schema DDL
--doc      produce documentation describing the schema
--mermaid  produce an ER diagram using Mermaid
```

The CLI is quiet on success — it reports the files it writes, surfaces validation
warnings on stderr, and exits non-zero only on validation errors.

 ## Try these examples

The repo includes three contracts you can run:

- `contracts/example.yml` — a clean four-table orders schema (the happy path)
- `contracts/cycle_example.yml` — a self-referencing table that triggers a cycle warning
- `contracts/bad_example.yml` — deliberately broken, to show validation output

```bash
go run ./cmd/cgen --input contracts/example.yml
go run ./cmd/cgen --input contracts/badExample.yml   # see the errors
```

## Contract Format

Contracts are YAML files describing tables, columns, constraints, and relationships.
For a cookbook of common patterns — composite keys, foreign keys, many-to-many,
self-references — see **[Writing Contracts](docs/writing-contracts.md)**.

### Supported constraints

- `not_null` — column that cannot be null

### Primary Key and Unique

Primary keys and unique constraints are declared at the table level and support
composite keys:

```yaml
primary_key: [id]                    # single-column PK
primary_key: [order_id, product_id]  # composite PK

unique:
  - [email]                          # single-column unique
  - [first_name, last_name]          # composite unique
```

### Foreign keys

Single-column foreign keys are declared on the referencing column:

```yaml
- name: customer_id
  type: uuid
  references:
    table: customers
    column: id
```

### Inference rules

The tool applies a small set of inferences to keep contracts concise:

- Primary key columns are automatically set to NOT NULL
- Primary key columns are implicitly unique (no separate UNIQUE constraint emitted)

Inferences are applied silently for now.

## Validation

ContractGen validates contracts in tiers, short-circuiting at the first tier that
fails so errors are reported at the right level of abstraction:

1. **Well-formedness** — contract structure and valid identifiers
2. **Schema validity** — completeness, duplicate tables/columns, foreign key and
   constraint resolution
3. **Generator concerns** — foreign key cycle detection (Kahn's algorithm)

Issues carry a severity. Errors block generation; warnings (such as FK cycles) are
surfaced but allow generation to continue.

## Features

- YAML contract parsing with multi-table support
- SQL DDL generation (Postgres dialect) with transactional wrapping
- Foreign key constraints, primary keys, composite keys, and unique constraints
- Markdown data dictionary generation
- Mermaid ER diagram generation
- JSON output for cross-tool interoperability
- Tiered schema validator with severity levels (errors vs warnings)
- Topological sort with cycle detection (Kahn's algorithm)
- CLI with selective output flags and configurable output directory
- Web UI with live generation across all output formats
- Test suite covering the validator and topological sort

## Planned Features

- [ ] Unify generators on the Schema layer (currently they consume Contract directly)
- [ ] Enrich JSON output with topological order and cycle metadata
- [ ] Multi-dialect support (currently just Postgres)
- [ ] Schema diffing between contract versions
- [ ] Reserved word and identifier escaping

## Architecture

```text
cmd/
├── cgen/         CLI entry point
└── cgen-web/     Web server entry point (HTMX + server-rendered templates)
internal/
├── parser/       YAML → Contract structs
├── normalise/    Applies inference rules to the parsed contract
├── validator/    Contract → tiered validation with severity
├── schema/       Contract → logical schema (ordered iteration, topological sort)
├── generators/   Contract → output formats (DDL, doc, JSON, mermaid)
├── web/          HTTP handlers, templates, static assets (embedded via go:embed)
└── cli/          Flag parsing and config
```

## Why This Exists

Data contracts are typically maintained manually across schemas, docs, and validation rules — leading to drift between what the contract says and what the database actually enforces.

ContractGen treats the YAML contract as the single source of truth and acts as a *compile target*: it parses, validates, and produces a structured Schema that downstream consumers (DDL, docs, diagrams, JSON) can render. The architecture follows the pattern used by tools like protoc, OpenAPI, and dbt — a typed specification with independent consumers.

## Deployment

The web UI ships as a single static Go binary in a multi-stage Docker build
(distroless final image) and is deployed on Fly.io. Templates and static
assets are embedded into the binary via `go:embed`, so the container has no
external file dependencies.

## Related Projects

- [SQL DQ Framework](https://github.com/Obiuno/metadata-sql-dq-framework) — validate data quality with SQL rules
- [DQ Synthetic Data](https://github.com/Obiuno/dq-synthetic-data) — generate synthetic data with controlled error injections

## License

MIT