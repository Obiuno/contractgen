# %%
"""
needs a lot of refactoring and adjusting to filepaths and a proper solution for extracting and using filepaths
"""
import yaml
from pprint import pprint

# YAML file path
yaml_file = '../contracts/example.yml'

# reading yaml data from a file
with open(yaml_file, 'r') as f:
    schema = yaml.safe_load(f)

pprint(schema['table'])

