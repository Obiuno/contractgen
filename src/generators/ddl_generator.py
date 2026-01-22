def generate_ddl(contract): 
        
        # define a full ddl to returned for all the tables
        full_ddl = []

        for tab in contract.tables:

              ddl_script = ""

              # add create and table name
              ddl_script +=  f"CREATE TABLE {tab['name'].upper()} (\n"
              
              # loop through columns to produce them
              # columns are a list of dictionaries
              
              # column for the names to be stored
              columns = []
              for col in tab['columns']:
                     col_str = f"\t{col['name']} {col['type'].upper()}"
                     
                     # primary key added
                     if col.get('primary_key'):
                            col_str += " PRIMARY KEY"

                     # unique constraint added
                     if col.get('unique'):
                            col_str += " UNIQUE"
                     
                     # add all the parts together
                     columns.append(col_str)

              ddl_script += ",\n".join(columns)
              ddl_script += "\n);"
              
              full_ddl.append(ddl_script)

        return "\n\n".join(full_ddl).strip()
# %%