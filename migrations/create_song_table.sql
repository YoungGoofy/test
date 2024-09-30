CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(100),
    song_title VARCHAR(100),
    release_date DATE,
    text TEXT,
    link VARCHAR(255)
);
