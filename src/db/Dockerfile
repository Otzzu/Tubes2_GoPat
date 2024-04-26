FROM postgres:latest

# Copy the initialization script and the dump file into the container
COPY init-s.sh /docker-entrypoint-initdb.d/init-script.sh
COPY ./dumps/dumpfile.sql /dumps/dumpfile.sql

# Make sure the script is executable
RUN chmod +x /docker-entrypoint-initdb.d/init-script.sh
