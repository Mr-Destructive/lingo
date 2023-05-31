CREATE TABLE IF NOT EXISTS user (
    id INT PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    name VARCHAR(100),
    email VARCHAR(50),
    password VARBINARY(255)
);

INSERT INTO user (id, username, name, email, password)
SELECT 1, 'test', 'test', 'test@lingo.com', 'test123'
WHERE NOT EXISTS (SELECT 1 FROM user WHERE id = 1);

INSERT INTO user (id, username, name, email, password)
SELECT 2, 'test2', 'test test', 'test2@lingo.com', 'test123'
WHERE NOT EXISTS (SELECT 1 FROM user WHERE id = 2);
