# %%
class Contract:
    pass

class SchemaContract (Contract):
    def __init__(self, data):
        self.table_name = data['table']['name']
        self.columns = data['table']['columns']
        self.actions = data['table']['actions']