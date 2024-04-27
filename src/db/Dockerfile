FROM postgres:latest

# Install dos2unix to convert line endings
# RUN apt-get update && apt-get install -y dos2unix

# # Copy the initialization script and the dump file into the container
# COPY init-s.sh /docker-entrypoint-initdb.d/init-script.sh
# COPY ./dumps/dumpfile.sql /dumps/dumpfile.sql

# # Convert line endings to Unix format and ensure the script is executable
# RUN dos2unix /docker-entrypoint-initdb.d/init-script.sh && \
#     chmod +x /docker-entrypoint-initdb.d/init-script.sh
