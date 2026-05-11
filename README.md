# ContractGen

Generate SQL DDL, JSON, and documentation from YAML data contracts.

⚠️ **v0.2 — Go rewrite, work in progress.** 

The original Python proof-of-concept has been ported to Go for a faster, single-binary tool. CLI and contract format are still evolving.

Define your schema once in YAML:

```yaml
version: "1.0"
description: "Core schema for the orders platform"

tables:
  - name: customers
    description: "Customer master data"
    columns:
      - name: id
        type: uuid
        constraints: ["primary_key", "not_null"]
      - name: email
        type: varchar
        constraints: ["not_null", "unique"]

  - name: orders
    description: "Customer orders"
    columns:
      - name: id
        type: uuid
        constraints: ["primary_key", "not_null"]
      - name: customer_id
        type: uuid
        constraints: ["not_null"]
        references:
          table: customers
          column: id
```

Generate Postgres DDL:

```bash
contractgen --input contracts/orders.yml --output schema.sql
```

Or output JSON for use by other tools:

```bash
contractgen --input contracts/orders.yml --json
```

## Installation

Requires Go 1.21+.

```bash
git clone https://github.com/Obiuno/contractgen
cd contractgen
go build -o contractgen ./cmd/cgen
```

Or run directly:

```bash
go run ./cmd/cgen --input contracts/example.yml
```

## Usage
```txt
contractgen --input <path> [--output <path>] [--json]
--input    Path to YAML contract file (required)
--output   SQL output path (default: schema.sql)
--json     Also print contract as JSON to stdout
```

## Contract Format

Contracts are YAML files describing tables, columns, constraints, and relationships.

### Supported column constraints
- `not_null` — column cannot be null
- `primary_key` — column is part of the primary key (currently single-column only)
- `unique` — column must be unique

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

Inferences are reported on stderr when generating output.

## Current Features (v0.2)
- ✅ YAML contract parsing (Go)
- ✅ Multi-table support
- ✅ SQL DDL generation with transactional wrapping
- ✅ Foreign key constraints
- ✅ JSON output for cross-language interop
- ✅ CLI with flag-based interface

## Planned Features
- [ ] Composite primary keys and unique constraints (table-level)
- [ ] Validation pass (FK targets exist, duplicate detection, type mismatches)
- [ ] Markdown documentation generation
- [ ] SELECT and INSERT template generation
- [ ] Multi-dialect support (currently Postgres-flavoured)
- [ ] Schema diffing between contract versions
- [ ] Reserved word and identifier escaping

## Architecture
```text
cmd/cgen/         CLI entry point
internal/
├── parser/       YAML → Contract structs
├── schema/       Contract → logical schema (lookup model)
├── generators/   Schema → output formats (DDL, JSON, etc.)
└── cli/          Flag parsing and config
```

## Status
**v0.1 - Initial public release.** Core functionality demonstrated, under active development.

## Why This Exists

Data contracts are foundational but typically maintained manually across schemas, docs, and validation rules — leading to drift between what the contract says and what the database actually enforces.

ContractGen treats the YAML contract as the single source of truth and generates the rest:
- DDL for database creation
- JSON for cross-tool interop
- (Planned) docs, validation rules, sample queries

The goal is a small, focused tool that does one thing well: keep your data contracts and generated artifacts in sync.

## Related Projects
- [SQL DQ Framework](https://github.com/Obiuno/metadata-sql-dq-framework) — validate data quality with SQL rules
- [DQ Synthetic Data](https://github.com/Obiuno/dq-synthetic-data) — generate synthetic data with controlled error injections

## License
MIT
