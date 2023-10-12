#!/bin/bash

# Set the database connection details
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="result"
DB_USER="postgres"
DB_PASSWORD="password"

# Set the path to the SQL migration file
SQL_FILE="./migration.sql"

# Run the migration
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $SQL_FILE
