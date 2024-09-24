#!/bin/bash

set -e

DB_NAME=${DB_NAME:mc_allowlist} 

echo "creating database $DB_NAME..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE $DB_NAME;
    GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $POSTGRES_USER;
EOSQL



echo "Creating table 'users' in database $DB_NAME..."   
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB_NAME" <<-EOSQL
    CREATE TABLE users (
        username VARCHAR(20) NOT NULL,
        message TEXT,
        request_date TIMESTAMP WITHOUT TIME ZONE,
        approval_date TIMESTAMP WITHOUT TIME ZONE,
        approved BOOLEAN,
        PRIMARY KEY (username)
    );
EOSQL


echo "Database and table created succesfully!"
