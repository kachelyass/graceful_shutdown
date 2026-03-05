CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);

INSERT INTO tasks (title)
VALUES
    ('Task number 1'),
    ('Task number 2');

