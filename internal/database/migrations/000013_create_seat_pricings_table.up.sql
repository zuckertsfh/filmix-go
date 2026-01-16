CREATE TABLE seat_pricings (
    id UUID NOT NULL UNIQUE,
    price BIGINT NOT NULL,
    day_type VARCHAR(255) NOT NULL,
    seat_type_id UUID NOT NULL,
    theater_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_seat_pricings_seat_type FOREIGN KEY (seat_type_id) REFERENCES seat_type(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_seat_pricings_theater FOREIGN KEY (theater_id) REFERENCES theaters(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);