CREATE TABLE IF NOT EXISTS routes (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name varchar(200) not null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc')
);