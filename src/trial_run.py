# %%
# first attempt
import yaml
from pprint import pprint 
from schema_contract import SchemaContract

with open('../contracts/example.yml', 'r') as f:
    data = yaml.safe_load(f)

pprint(data)


# %%
# separated generators

import yaml
from schema_contract import SchemaContract
from generators.doc_generator import generate_doc
from generators.ddl_generator import generate_ddl
from generators.select_generator import generate_select

with open('../contracts/example.yml', 'r') as f:
    data = yaml.safe_load(f)

ex_contract = SchemaContract(data)
print(generate_doc(ex_contract))
print(generate_ddl(ex_contract))
print(generate_select(ex_contract))
# %%
import yaml
from schema_contract import Contract, SchemaContract
from generators.ddl_generator import generate_ddl

ex_contract = Contract.from_yaml('../contracts/example.yml')


print("--- DDL Script ---")
print(generate_ddl(ex_contract))
# %%
