CREATE TABLE IF NOT EXISTS user (
    username VARCHAR(50) UNIQUE,
    name VARCHAR(100),
    password_hash VARBINARY(255)
);
