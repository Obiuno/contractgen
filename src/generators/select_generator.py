# %%
def generate_select(contract):
        select_script = ""

        # add select and table name
        select_script +=  f"SELECT\n"
        
        # loop through columns to produce them
        # columns are a list of dictionaries
        
        # column for the names to be stored
        columns = []
        for col in contract.columns:
               col_str = f"\t{col['name']}"
               
               # add all the parts together
               columns.append(col_str)

        select_script += ",\n".join(columns)

        # add the from and table
        select_script += f"\nFROM {contract.table_name};"

        return select_script

# %%
