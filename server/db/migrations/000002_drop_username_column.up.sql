CREATE TABLE IF NOT EXISTS users_new (
    id           TEXT NOT NULL PRIMARY KEY,
    email        TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

INSERT INTO users_new (id, email, password_hash)
SELECT id, email, password_hash FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;
