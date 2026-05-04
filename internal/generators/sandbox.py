# %%
import yaml
from pprint import pprint

# YAML file path
yaml_file = '../contracts/example.yml'

# reading yaml data from a file
with open(yaml_file, 'r') as f:
    contract = yaml.safe_load(f)

pprint(contract['table'])
# %%
def format_column(column):
      col_name = column['name']
      col_type = column['type']
      col_desc = column['description']
      return f"{col_name}   ({col_type})"


# %%
# generating a document

#function for generating tables
#make table name
#table decription
#fields
#list out actions - need to make an optional input

#table → dict

#table['name'] → string

#table['columns'] → list

#table['actions'] → list (optional)

def generate_doc(table):
        table_name = table['name']
        print(f'Table name:\n{table_name.upper()}\n')

        column = table['columns']
        
        print('Columns:')
        for col in column:
              column_name = col['name']
              column_type = col['type']
              print(f'{column_name} ({column_type})')

generate_doc(contract['table'])


# %%
