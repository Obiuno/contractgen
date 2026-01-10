# %%
def generate_doc(contract):
        schema_doc = ""

        # produce table name
        schema_doc +=  f'Table name:\n{contract.table_name.upper()}\n\n'
        

        # column name and type
        schema_doc += 'Columns:\n'
        # creating table view for columns
        schema_doc += "| Column | Type | Description |\n"
        schema_doc += "|--------|------|-------------|\n"
        for col in contract.columns:
              column_name = col['name']
              column_type = col['type']
              column_desc = col.get('description', '')
              schema_doc += f'| {column_name} | {column_type} | {column_desc} | \n'
        
        return schema_doc
# %%
