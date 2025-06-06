import mysql.connector
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

def get_db_connection():
    """
    Establishes a connection to the MySQL database using environment variables.

    Returns:
        mysql.connector.MySQLConnection: A MySQL connection object.
    Raises:
        mysql.connector.Error: If the connection to the database fails.
    """
    try:
        connection = mysql.connector.connect(
            host=os.getenv("DB_HOST"),
            user=os.getenv("DB_USER"),
            password=os.getenv("DB_PASSWORD"),
            database=os.getenv("DB_NAME")
        )
        return connection
    except mysql.connector.Error as err:
        print(f"Error connecting to database: {err}")
        raise

def close_db_connection(connection):
    """
    Closes the provided database connection.

    Args:
        connection (mysql.connector.MySQLConnection): The MySQL connection object to close.
    """
    if connection and connection.is_connected():
        connection.close()
        print("Database connection closed.")

def create_db_cursor(connection):
    """
    Creates a cursor object for the given database connection.

    Args:
        connection (mysql.connector.MySQLConnection): The MySQL connection object.

    Returns:
        mysql.connector.cursor.MySQLCursor: A cursor object for executing queries.
    """
    return connection.cursor()

def execute_query(cursor, query, params=None):
    """
    Executes a given SQL query with optional parameters.

    Args:
        cursor (mysql.connector.cursor.MySQLCursor): The database cursor.
        query (str): The SQL query to execute.
        params (tuple, optional): Parameters to pass to the query. Defaults to None.

    Returns:
        None
    Raises:
        mysql.connector.Error: If the query execution fails.
    """
    try:
        if params:
            cursor.execute(query, params)
        else:
            cursor.execute(query)
    except mysql.connector.Error as err:
        print(f"Error executing query: {err}")
        raise

def fetch_all(cursor):
    """
    Fetches all results from a query.

    Args:
        cursor (mysql.connector.cursor.MySQLCursor): The database cursor.

    Returns:
        list: A list of tuples representing the query results.
    """
    return cursor.fetchall()

def fetch_one(cursor):
    """
    Fetches a single result from a query.

    Args:
        cursor (mysql.connector.cursor.MySQLCursor): The database cursor.

    Returns:
        tuple: A tuple representing the query result, or None if no result is found.
    """
    return cursor.fetchone()

def commit_changes(connection):
    """
    Commits the changes made to the database.

    Args:
        connection (mysql.connector.MySQLConnection): The MySQL connection object.
    """
    connection.commit()