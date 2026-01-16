CREATE TABLE movie_genres (
    id UUID NOT NULL UNIQUE,
    genre VARCHAR(255) NOT NULL,
    PRIMARY KEY(id)
);