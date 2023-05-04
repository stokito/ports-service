# Ports Service
A Port domain API service to store and retrieve ports.

Each port looks in JSON like: 
```json
{
  "name": "Djibouti",
  "city": "Djibouti",
  "country": "Djibouti",
  "alias": [],
  "regions": [],
  "coordinates": [
    43.1456475,
    11.5720765
  ],
  "province": "Djibouti",
  "timezone": "Africa/Djibouti",
  "unlocs": [
    "DJJIB",
    "DJPOD"
  ],
  "code": "77701"
}
```

So the port may have many `"unlocs": ["DJJIB", "DJPOD"]` and we need a deduplication by the unloc code.

The repo contains:
* ports-service - a REST-like API service to save and get a port.
* ports-import tool to load ports from JSON file e.g. `ports.json`.

They both can be configured with a config.json file e.g. DB settings, API listen address etc.

## API
Sample of API. Create/update a port with POST:
```
POST http://localhost:8080/api/v1/ports/
Authorization: Basic api secret
Content-Type: application/json

{
  "name": "Ajman",
  "unlocs": [
    "AEAJM"
  ]
}
```

Get/retrieve with GET:
```
GET http://localhost:8080/api/v1/ports/?unloc=AEAJM
Authorization: Basic api secret
```

Currently, for the API only HTTP is supported so use a reverse proxy.
You can configure a basic authorization for the API.

## Setup
The service use PostgreSQL with JSON columns support.
It's database can be initialized with `init.sql` file.

The provided docker-compose uses an official [PostgreSQL image](https://hub.docker.com/_/postgres) with the init.sql script.

Use `docker-compose up` and `docker-compose down` to start and stop the application.

Default configuration is:
```json
{
  "DatabaseUrl": "postgres://postgres:postgres@db:5432/portsdb?search_path=ports_schema",
  "PortsFilePath": "ports.json",
  "ListenAddr": ":8080",
  "Credentials": {
    "api": "secret"
  }
}
```

## Local development
You can override config options like DB URL in the `config.local.json`.
For Goland run configuration were added.
The `requests.http` file contains samples of API requests.

## Build and Usage
You can build manually with:

    go build  -o ports-import ./cmd/import
    go build  -o ports-service ./cmd/service

Or you can run `make build`.
To run a linter execute `make run-lint`.
You may need to install a linter `sudo snap install golangci-lint`.

## Technical test

- Given a file with ports data (ports.json), write a port domain service that either creates a new record in a database, or updates the existing one (Hint: no need for delete or other methods).
- The file is of unknown size, it can contain several millions of records, you will not be able to read the entire file at once.
- The service has limited resources available (e.g. 200MB ram).
- The end result should be a database containing the ports, representing the latest version found in the JSON. (Hint: use an in memory database to save time and avoid complexity).
- A Dockerfile should be used to contain and run the service (Hint: extra points for avoiding compilation in docker).
- Provide at least one example per test type that you think are needed for your assignment. This will allow the reviewer to evaluate your critical thinking as well as your knowledge about testing.
- Your readme.md should explain how to run your program and test it.
- The service should handle certain signals correctly (e.g. a TERM or KILL signal should result in a graceful shutdown).

### Bonus points

- Address security concerns for Docker
- Database in docker container
- Docker-compose file
