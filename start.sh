#!/bin/sh

# set dash e: we use set -e instruction to make sure that the script will exit immediately if a command return a non-zero status
set -e

echo "run db migration"

# use -verbose option to print out all details when the migration is run
# finally the up argument is used to run migrate up all migrations 
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# after we run migrate up we start the app
echo "start the app"

# Basically means: takes all parameters passed to the script and run it
exec "$@"