CREATE TABLE IF NOT EXISTS user (
    id INT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARBINARY(255) NOT NULL
);

INSERT INTO user (id, username, email, password)
SELECT 1, 'test', 'test@lingo.com', 'test123'
WHERE NOT EXISTS (SELECT 1 FROM user WHERE id = 1);

INSERT INTO user (id, username, email, password)
SELECT 2, 'test2', 'test2@lingo.com', 'test123'
WHERE NOT EXISTS (SELECT 1 FROM user WHERE id = 2);
