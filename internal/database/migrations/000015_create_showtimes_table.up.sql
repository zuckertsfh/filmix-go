CREATE TABLE showtimes (
    id UUID NOT NULL UNIQUE,
    status BOOLEAN NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    movie_id UUID NOT NULL,
    studio_id UUID NOT NULL,
    theater_id UUID NOT NULL,
    seat_pricing_id UUID NOT NULL,
    seat_pricing_override_id UUID,
    PRIMARY KEY(id),
    CONSTRAINT fk_showtimes_movie FOREIGN KEY (movie_id) REFERENCES movies(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_showtimes_studio FOREIGN KEY (studio_id) REFERENCES studios(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_showtimes_theater FOREIGN KEY (theater_id) REFERENCES theaters(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_showtimes_seat_pricing FOREIGN KEY (seat_pricing_id) REFERENCES seat_pricings(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_showtimes_override FOREIGN KEY (seat_pricing_override_id) REFERENCES seat_pricing_overrides(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);