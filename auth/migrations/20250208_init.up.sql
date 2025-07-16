CREATE TABLE IF NOT EXISTS people (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL
);

INSERT INTO people (username) VALUES ('chel');
INSERT INTO people (username) VALUES ('chelick');
