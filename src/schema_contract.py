# %%
import yaml

class Contract:
    @staticmethod
    def from_yaml(file_path):
        with open(file_path, 'r') as f:
            data = yaml.safe_load(f)

        if 'tables' in data:
            return SchemaContract(data)
        elif 'process' in data:
            pass
        else:
            raise ValueError("Invalid Contract: must define 'tables' or 'process'")

class SchemaContract(Contract):
    def __init__(self, data):
        # need logic for single table contract and multi table contract
        # when single table it is a dict, when multiple tables it is a list of dics
        tables_data = data['tables']

        if isinstance(tables_data, dict):
            tables_data = [tables_data]
        
        self.tables = tables_data
# %%
