CREATE TABLE movies (
    id UUID NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    tagline VARCHAR(255) NOT NULL,
    overview TEXT NOT NULL,
    poster_url TEXT NOT NULL,
    backdrop_url TEXT NOT NULL,
    trailer_url TEXT NOT NULL,
    duration INTEGER NOT NULL,
    popularity INTEGER NOT NULL,
    movie_status_id UUID NOT NULL,
    movie_rating_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_movies_status FOREIGN KEY (movie_status_id) REFERENCES movie_statuses(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_movies_rating FOREIGN KEY (movie_rating_id) REFERENCES movie_ratings(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);