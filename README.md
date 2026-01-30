## ContracGen

Generate SQL DDL and documentation from YAML data contracts.

⚠️ **v0.1 - Initial Release**: Core concept working 
This is an early proof-of-concept.

## What It Does

Define your schema in YAML:
```yaml
tables:
  name: users
  columns:
    - name: id
      type: integer
      primary_key: true
      description: User ID
```

Generate SQL DDL:
```python
from schema_contract import Contract
from generators.ddl_generator import generate_ddl

contract = Contract.from_yaml('contracts/example.yml')
print(generate_ddl(contract))
```

## Current Features (v0.1)
- ✅ YAML contract parsing
- ✅ Single and multi-table support
- ✅ Basic SQL DDL generation
- ✅ Markdown documentation generation
- ✅ SELECT statement generation

## Planned Features (v0.2+)
- [ ] SQL identifier validation and escaping
- [ ] Warnings for reserved words and special characters
- [ ] Foreign key relationships
- [ ] Data dictionary export (CSV/JSON)
- [ ] Comprehensive tests
- [ ] process visualisation

## Status
**v0.1 - Initial public release.** Core functionality demonstrated, under active development.

## Why This Exists
Data contracts are critical but often documented manually. ContractGen aims to:
- Serve as the Single Source of the Truth
- Define schemas once (YAML)
- Generate everything else (SQL, docs, dictionaries)
- Enforce best practices through validation
- Enable contract-driven development

## Installation
```bash
git clone https://github.com/Obiuno/contractgen
cd contractgen
# Run examples in trial_run.py
```

## Realted Projects

- [SQL DQ Framework](https://github.com/Obiuno/metadata-sql-dq-framework) - Validate data quality with SQL rules
- [DQ Synthetic Data](https://github.com/Obiuno/dq-synthetic-data) - Generate test data for validation

## Contributing
This is v0.1 - feedback and contributions welcome as I develop this further.

## License
MIT
