CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INT NOT NULL,
    color VARCHAR(16),
    avatar VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50),
    url VARCHAR(255) UNIQUE,
    user_id INT,
    profile_id INT,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE ON UPDATE CASCADE
    CONSTRAINT uq_name_user UNIQUE (name, user_id)
);

INSERT INTO profiles(id, user_id, color, avatar)
SELECT 1, 1, '#decade', 'https://www.gravatar.com/avatar/00000000000000000000000000000000?d=mp&f=y'
WHERE NOT EXISTS (SELECT 1 FROM profiles WHERE id = 1);

INSERT INTO links(id, name, url, user_id, profile_id)
SELECT 1, 'twitter', 'twitter.com/test', 1, 1
WHERE NOT EXISTS (SELECT 1 FROM links WHERE id = 1);

INSERT INTO links(id, name, url, user_id)
SELECT 2, 'github', 'https://github.com/test', 2
WHERE NOT EXISTS (SELECT 1 FROM links WHERE id = 2);

INSERT INTO links(id, name, url, user_id)
SELECT 3, 'twitter', 'https://twitter.com/test2', 2
WHERE NOT EXISTS (SELECT 1 FROM links WHERE id = 3);

INSERT INTO links(id, name, url, user_id)
SELECT 4, 'github', 'https://github.com/test1', 1
WHERE NOT EXISTS (SELECT 1 FROM links WHERE id = 4);
