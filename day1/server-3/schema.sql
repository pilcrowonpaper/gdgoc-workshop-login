CREATE TABLE user (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash BLOB NOT NULL,
    password_salt BLOB NOT NULL
) STRICT;