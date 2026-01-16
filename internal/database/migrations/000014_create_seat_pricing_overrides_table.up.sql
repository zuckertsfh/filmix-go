CREATE TABLE seat_pricing_overrides (
    id UUID NOT NULL UNIQUE,
    price BIGINT NOT NULL,
    notes TEXT NOT NULL,
    movie_id UUID NOT NULL,
    seat_type_id UUID NOT NULL,
    theater_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_overrides_movie FOREIGN KEY (movie_id) REFERENCES movies(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_overrides_seat_type FOREIGN KEY (seat_type_id) REFERENCES seat_type(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_overrides_theater FOREIGN KEY (theater_id) REFERENCES theaters(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);