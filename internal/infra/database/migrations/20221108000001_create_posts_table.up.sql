CREATE TABLE IF NOT EXISTS nix_education.posts (
    id              serial PRIMARY KEY,
    user_id           int NOT NULL,
    title        varchar(100) NOT NULL,
    body          varchar(1000) NOT NULL
);
