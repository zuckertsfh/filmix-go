CREATE TABLE movie_ratings (
    id UUID NOT NULL UNIQUE,
    rating VARCHAR(255) NOT NULL,
    PRIMARY KEY(id)
);