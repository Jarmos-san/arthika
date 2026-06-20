CREATE TABLE IF NOT EXISTS users_old (
    id           TEXT NOT NULL PRIMARY KEY,
    username     TEXT NOT NULL UNIQUE,
    email        TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

INSERT INTO users_old (id, username, email, password_hash)
SELECT id, '', email, password_hash FROM users;

DROP TABLE users;

ALTER TABLE users_old RENAME TO users;
