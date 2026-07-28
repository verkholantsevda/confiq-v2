CREATE TABLE users_new (
    id INTEGER NOT NULL,
    username VARCHAR(80) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    config_limit INTEGER,
    is_admin BOOLEAN,
    group_id INTEGER,
    created_at DATETIME,
    PRIMARY KEY (id),
    UNIQUE(username),
    FOREIGN KEY(group_id) REFERENCES groups(id)
);

INSERT INTO users_new (
    id,
    username,
    password_hash,
    config_limit,
    is_admin,
    group_id,
    created_at
)
SELECT
    id,
    username,
    password_hash,
    config_limit,
    is_admin,
    group_id,
    created_at
FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;