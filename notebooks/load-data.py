#----------------------------------------------------------------------------------------
# import packages
#----------------------------------------------------------------------------------------
from dotenv import load_dotenv
import os
import sqlite3
import pandas as pd
import psycopg2
import psycopg2.extras
from datetime import datetime
import logging

#----------------------------------------------------------------------------------------
# Set up logging
#----------------------------------------------------------------------------------------
date_format = "%Y-%m-%d %H:%M:%S"
formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s', datefmt=date_format)

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)
handler = logging.StreamHandler()
handler.setFormatter(formatter)
logger.addHandler(handler)

#----------------------------------------------------------------------------------------
# load environment variables
#----------------------------------------------------------------------------------------
load_dotenv()
loading_dir:str|None = os.getenv('LOADING_DIR')
db_host:str|None = os.getenv('DB_HOST')
db_port:str|None = os.getenv('DB_PORT')
db:str|None = os.getenv('DB')
db_user:str|None = os.getenv('DB_USER')
db_pw:str|None = os.getenv('DB_PW')
print(f"db_host:{db_host}, db_port:{db_port}, db:{db}, db_user:{db_user}, db_pw:{db_pw}")

#----------------------------------------------------------------------------------------
# setup vars and db connection
#----------------------------------------------------------------------------------------
db_files:list[str] = os.listdir(loading_dir)
psql_conn = psycopg2.connect( host=db_host, database=db, user=db_user, password=db_pw, port=db_port)
psql_conn.autocommit = True
pg_cur = psql_conn.cursor()
print(pg_cur.execute("SELECT CURRENT_USER;"))


#----------------------------------------------------------------------------------------
# define and call load method
#----------------------------------------------------------------------------------------
def load_db_file(filename:str):
    logger.info(f"running load_db_file('{filename}')")
    # Connect to the SQLite database
    conn = sqlite3.connect(os.path.join(str(loading_dir), filename))

    # Create a cursor object
    sql_cu = conn.cursor()

    # List all tables in the database
    sql_cu.execute("SELECT name, sql FROM sqlite_master WHERE type='table'")

    # Fetch all rows returned by the query
    tables = sql_cu.fetchall()

    # hold table names
    table_names:list = []
    for table in tables:
        table_names.append(table[0])

    # process for each table
    for table in tables:
        logger.info(f"loading table {table[0]}")
        if table[0] not in ['sqlite_sequence']:
            #----------------------------------------------------------------------------------------
            # staging
            #----------------------------------------------------------------------------------------
            # drop table if exists
            sql: str = f"drop table if exists staging.{table[0]} cascade"
            logger.info(f"  dropping")
            pg_cur.execute(sql)

            # modify create table statement
            sql = table[1].replace(f"CREATE TABLE {table[0]} (", f"CREATE TABLE staging.{table[0]} (")
            sql = sql.replace("id INTEGER PRIMARY KEY AUTOINCREMENT,", "id SERIAL PRIMARY KEY,")

            for tn in table_names:
                sql = sql.replace(f"REFERENCES {tn}(", f"REFERENCES staging.{tn}(")

            logger.info(f"  creating staging")
            pg_cur.execute(sql)

            copy_data_to_staging(sql_cu, pg_cur, table[0])
            #----------------------------------------------------------------------------------------
            # bronze
            #----------------------------------------------------------------------------------------
            if not table_exists(pg_cur, table[0]):
            
                # modify create table statement
                sql = table[1].replace(f"CREATE TABLE {table[0]} (", f"CREATE TABLE bronze.{table[0]} (")
                sql = sql.replace("id INTEGER PRIMARY KEY AUTOINCREMENT,", "id SERIAL PRIMARY KEY,")

                for tn in table_names:
                    sql = sql.replace(f"REFERENCES {tn}(", f"REFERENCES bronze.{tn}(")

                logger.info(f" creating bronze")
                pg_cur.execute(sql)
        else:
            logger.info(f"  skipped")


def table_exists(cursor, table_name) -> bool:
    try:
        cursor.execute(f"SELECT 1 FROM information_schema.tables WHERE table_name = '{table_name}' and table_schema = 'bronze'")
        result = cursor.fetchone()
        return bool(result)
    except psycopg2.Error:
        return False

def copy_data_to_staging(sql_cu, pg_cur, table_name):

    sql_cu.execute(f"PRAGMA table_info({table_name})")
    columns = [row[1] for row in sql_cu.fetchall() if row[1] != 'id']
    print(list_to_string(columns))

    sql_cu.execute(f"SELECT * FROM {table_name}")
    rows = sql_cu.fetchall()

def list_to_string(l):
    return ','.join([f"'{item}'" for item in l])

load_db_file("journeys_nzta_20240907100303.db")