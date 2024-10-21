-- +migrate Up notransaction
CREATE TABLE employees (
    id BIGSERIAL NOT NULL,
    name text NOT NULL,
    position text NOT NULL,
    salary float8 NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT employees_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE employees;