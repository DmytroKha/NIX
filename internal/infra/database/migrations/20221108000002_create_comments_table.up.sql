CREATE TABLE IF NOT EXISTS nix_education.comments (
    id              serial PRIMARY KEY,
    post_id           int NOT NULL,
    name         varchar(100) NOT NULL,
    email        varchar(100) NOT NULL,
    body          varchar(1000) NOT NULL
);
