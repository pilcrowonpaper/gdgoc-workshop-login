CREATE TABLE user (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash BLOB NOT NULL,
    password_salt BLOB NOT NULL
) STRICT;

CREATE TABLE session (
    id TEXT NOT NULL PRIMARY KEY,
    secret_hash BLOB NOT NULL,
    user_id TEXT NOT NULL REFERENCES user(id),
    created_at INTEGER NOT NULL
) STRICT;