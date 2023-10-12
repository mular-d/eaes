@echo off

REM Set the database connection details
set DB_HOST=localhost
set DB_PORT=5432
set DB_NAME=result
set DB_USER=postgres
set DB_PASSWORD=password

REM Set the path to the SQL migration file
set SQL_FILE=%~dp0migration.sql

REM Run the migration
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -f %SQL_FILE%
