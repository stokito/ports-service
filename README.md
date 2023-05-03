# Ports Service

## Setup
The service use PostgreSQL with JSON columns support.
It's database can be initialized with `init.sql` file.

The provided docker-compose uses an official [PostgreSQL image](https://hub.docker.com/_/postgres) with the init.sql script.

Use `docker-compose up` and `docker-compose down` to start and stop the application.

## Local development
You can override config options like DB URL in the `config.local.json`.
The `requests.http` file contains samples of API requests.
