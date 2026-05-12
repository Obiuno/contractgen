# ContractGen

Generate SQL DDL, JSON, and documentation from YAML data contracts.

⚠️ **v0.2 — Go rewrite**

Rewrite completed for core generators, validator and CLI completed in Go.
Testing and additional features in progress.

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
if no flags specified, defaults to produce all 3

--input    Path to YAML contract file (required)
--output-dir   directory where files are generated

Boolean flags
--json     produces schema in JSON format
--ddl      produces schema DDL
--doc      produces doc describing schema
```

## Contract Format

Contracts are YAML files describing tables, columns, constraints, and relationships.

### Supported constraints

- `not_null` — column that cannot be null

### Primary Key and Unique
Primary Keys and Unique fields are declared at the table level and support composite keys

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

## Current Features (v0.2)

- ✅ YAML contract parsing (Go)
- ✅ Multi-table support
- ✅ SQL DDL generation with transactional wrapping
- ✅ Foreign key constraints
- ✅ JSON output for cross-language usage
- ✅ CLI with flag-based interface for different outputs
- ✅ Composite primary keys and unique constraints (table-level)
- ✅ Validator for: duplicate tables, duplicate columns and foreign keys
- ✅ Markdown document generation

## Planned Features

- [ ] Multi-dialect support (currently Postgres-flavoured)
- [ ] Schema diffing between contract versions
- [ ] Reserved word and identifier escaping

## Architecture

```text
cmd/cgen/         CLI entry point
internal/
├── parser/       YAML → Contract structs
├── validator/    Contract → validation
├── schema/       Contract → logical schema (lookup model)
├── generators/   Schema → output formats (DDL, JSON, Doc)
└── cli/          Flag parsing and config
```

## Why This Exists

Data contracts are foundational but typically maintained manually across schemas, docs, and validation rules — leading to drift between what the contract says and what the database actually enforces.

ContractGen treats the YAML contract as the single source of truth and generates the rest:

- DDL for database creation
- JSON for cross-tool interop
- Docs as a starting point for a data dictionary
- (Planned) validation rules, sample queries

The goal is a small, focused tool that does one thing well: keep your data contracts and generated artifacts in sync.

## Related Projects

- [SQL DQ Framework](https://github.com/Obiuno/metadata-sql-dq-framework) — validate data quality with SQL rules
- [DQ Synthetic Data](https://github.com/Obiuno/dq-synthetic-data) — generate synthetic data with controlled error injections

## License

MIT
