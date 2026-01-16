CREATE TABLE genre_movie (
    id UUID NOT NULL UNIQUE,
    movie_id UUID NOT NULL,
    movie_genre_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_genre_movie_movie FOREIGN KEY (movie_id) REFERENCES movies(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_genre_movie_genre FOREIGN KEY (movie_genre_id) REFERENCES movie_genres(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);