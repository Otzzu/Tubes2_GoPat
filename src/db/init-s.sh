#!/bin/bash
# Initialization script for PostgreSQL Docker container

set -e

# Display a log message
echo "Running PostgreSQL initialization script..."

# Ensure environment variables are available
DB=${POSTGRES_DB:-postgres}
USER=${POSTGRES_USER:-postgres}

# Location of the dump file to be loaded
DUMP_FILE="/dumps/dumpfile.sql"

# Check if the dump file exists
if [ -f "$DUMP_FILE" ]; then
    echo "Dump file found. Starting to load the dump file..."
    psql -U "$USER" -d "$DB" -f "$DUMP_FILE"
    echo "Dump file loaded successfully."
else
    echo "Dump file not found: $DUMP_FILE"
fi

echo "Initialization script completed successfully."
