CREATE TABLE guestbook (
    id serial PRIMARY KEY,
    name VARCHAR(50),
    is_approved BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE fitness_recaps (
    id serial PRIMARY KEY,
    weight INTEGER,
    distance INTEGER,
    date TIMESTAMP
);