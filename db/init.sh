#!/bin/bash
set -e

#
# $POSTGRES_PASSWORD
#

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER rdev;
	CREATE DATABASE backoffice_db;
	GRANT ALL PRIVILEGES ON DATABASE backoffice_db TO rdev;
EOSQL