-- +migrate Up
-- +migrate StatementBegin

-- Tabel Animal
CREATE TABLE Animals (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    class VARCHAR(100) NOT NULL,
    legs SMALLINT NOT NULL
);

-- +migrate StatementEnd
