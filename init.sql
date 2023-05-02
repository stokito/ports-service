CREATE SCHEMA ports_schema;
SET search_path TO ports_schema;

CREATE TABLE ports_schema.ports
(
    unlocks VARCHAR NOT NULL
        PRIMARY KEY
        UNIQUE,
    port    JSON    NOT NULL
);
