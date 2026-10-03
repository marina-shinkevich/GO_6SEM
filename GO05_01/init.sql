CREATE TABLE IF NOT EXISTS celebrities (
    id            SERIAL PRIMARY KEY,
    full_name     VARCHAR(255) NOT NULL,
    nationality   VARCHAR(100) NOT NULL,
    req_photo_path VARCHAR(255) NOT NULL
);

INSERT INTO celebrities (full_name, nationality, req_photo_path) VALUES
    ('Leonardo DiCaprio', 'American',  '/photos/dicaprio.jpg'),
    ('Scarlett Johansson', 'American', '/photos/johansson.jpg'),
    ('Christoph Waltz',   'Austrian',  '/photos/waltz.jpg');