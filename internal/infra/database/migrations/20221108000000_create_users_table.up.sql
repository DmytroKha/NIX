CREATE TABLE IF NOT EXISTS nix_education.users (
    id              serial PRIMARY KEY,
    email           varchar(100) NOT NULL,
    password        varchar(100) NOT NULL,
    name         varchar(100) NOT NULL
);
