# Writing Contracts

A cookbook of common patterns for ContractGen YAML contracts. Each recipe shows a
self-contained snippet you can adapt. For the full field reference, see the
[README](../README.md#contract-format).

## Contents

- [The minimal contract](#the-minimal-contract)
- [A single table](#a-single-table)
- [Single-column primary key](#single-column-primary-key)
- [Composite primary key](#composite-primary-key)
- [Unique constraints](#unique-constraints)
- [Foreign keys](#foreign-keys)
- [Many-to-many (join tables)](#many-to-many-join-tables)
- [Self-referencing tables](#self-referencing-tables)
- [Avoiding foreign key cycles](#avoiding-foreign-key-cycles)
- [What the inference rules do for you](#what-the-inference-rules-do-for-you)
- [Common mistakes](#common-mistakes)

## The minimal contract

Every contract needs a `version` and at least one table with one column:

```yaml
version: "1.0"

tables:
  - name: customers
    columns:
      - name: id
        type: uuid
        constraints: ["not_null"]
    primary_key: [id]
```

`description` is optional at both the contract and table level, but worth adding —
it flows through to the generated documentation.

## A single table

A fuller table with a description and several columns:

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
```

## Single-column primary key

Declare the primary key at the **table level**, not on the column. It takes a list
even when there's only one column:

```yaml
    columns:
      - name: id
        type: uuid
    primary_key: [id]
```

You don't need to add `not_null` to a primary key column — see
[inference rules](#what-the-inference-rules-do-for-you).

## Composite primary key

List multiple columns. Common for join tables and line-item tables:

```yaml
  - name: order_items
    columns:
      - name: order_id
        type: uuid
      - name: product_id
        type: uuid
      - name: quantity
        type: integer
        constraints: ["not_null"]
    primary_key: [order_id, product_id]
```

The combination of `order_id` + `product_id` must be unique; individual columns can
repeat.

## Unique constraints

Declared at the table level under `unique`, as a list of column groups. Each inner
list is one constraint:

```yaml
  - name: customers
    columns:
      - name: id
        type: uuid
      - name: email
        type: varchar
        constraints: ["not_null"]
      - name: first_name
        type: varchar
      - name: last_name
        type: varchar
    primary_key: [id]
    unique:
      - [email]                    # email alone must be unique
      - [first_name, last_name]    # the pair must be unique together
```

Note the difference: `[first_name, last_name]` means the *combination* is unique
(two people can share a first name, just not a full name). Two separate single-column
constraints would mean each is independently unique, which is usually not what you
want.

## Foreign keys

Declared on the **referencing column** via `references`:

```yaml
  - name: orders
    columns:
      - name: id
        type: uuid
      - name: customer_id
        type: uuid
        constraints: ["not_null"]
        references:
          table: customers
          column: id
    primary_key: [id]
```

The referenced table and column must exist in the same contract — ContractGen
validates this and will report an error if `customers.id` isn't defined.

Foreign keys are currently single-column only. Composite foreign keys are not yet
supported.

## Many-to-many (join tables)

There's no special syntax — a many-to-many is just a table with two foreign keys
and a composite primary key over them:

```yaml
  - name: students
    columns:
      - name: id
        type: uuid
    primary_key: [id]

  - name: courses
    columns:
      - name: id
        type: uuid
    primary_key: [id]

  - name: enrolments
    description: "Links students to the courses they take"
    columns:
      - name: student_id
        type: uuid
        constraints: ["not_null"]
        references:
          table: students
          column: id
      - name: course_id
        type: uuid
        constraints: ["not_null"]
        references:
          table: courses
          column: id
    primary_key: [student_id, course_id]
```

## Self-referencing tables

A table can reference itself — common for hierarchies (org charts, category trees).
The foreign key column should be nullable (no `not_null`) so the root row can have
no parent:

```yaml
  - name: employees
    columns:
      - name: id
        type: uuid
      - name: name
        type: varchar
        constraints: ["not_null"]
      - name: manager_id
        type: uuid
        references:
          table: employees
          column: id
    primary_key: [id]
```

A self-reference is technically a cycle of length one. ContractGen surfaces it as a
**warning**, not an error — generation still proceeds. See below.

## Avoiding foreign key cycles

ContractGen runs a topological sort (Kahn's algorithm) to order tables by their
dependencies. If table A references B and B references A, there's no valid creation
order, and you'll get a cycle warning:

```yaml
  - name: a
    columns:
      - name: id
        type: uuid
      - name: b_id
        type: uuid
        references:
          table: b
          column: id
    primary_key: [id]

  - name: b
    columns:
      - name: id
        type: uuid
      - name: a_id
        type: uuid
        references:
          table: a
          column: id
    primary_key: [id]
```

This is a *warning*, not an error — the generators still run. But it usually signals
a modelling problem. The common fix is to make one side nullable and populate it in a
second pass, or to introduce a join table that breaks the cycle.

## What the inference rules do for you

A couple of things are applied automatically so you don't have to write them:

- **Primary key columns become NOT NULL.** You don't need to add `not_null` to a PK
  column; it's implied.
- **Primary key columns are implicitly unique.** No separate UNIQUE constraint is
  emitted for the primary key.

So this:

```yaml
      - name: id
        type: uuid
    primary_key: [id]
```

is equivalent to writing `constraints: ["not_null"]` on `id` and adding a unique
constraint — ContractGen does it for you.

## Common mistakes

**Putting `primary_key` on the column instead of the table.** It's a table-level
field:

```yaml
# wrong
      - name: id
        type: uuid
        primary_key: true

# right
      - name: id
        type: uuid
    primary_key: [id]
```

**Forgetting that `primary_key` and `unique` take lists.** Even a single-column key
is a list: `primary_key: [id]`, not `primary_key: id`.

**Referencing a table or column that doesn't exist.** Foreign key targets must be
defined in the same contract. A typo in `references.table` or `references.column`
will fail validation.

**Confusing a composite unique constraint with two single-column ones.**
`unique: [[a, b]]` means the pair is unique; `unique: [[a], [b]]` means each is
independently unique. They're very different constraints.

**Adding redundant `not_null` to primary key columns.** Harmless, but unnecessary —
it's inferred.