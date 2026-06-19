DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id           TEXT NOT NULL PRIMARY KEY,
    username     TEXT NOT NULL UNIQUE,
    email        TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);
