CREATE TABLE seats (
    id UUID NOT NULL UNIQUE,
    row VARCHAR(255) NOT NULL,
    number INTEGER NOT NULL,
    active BOOLEAN NOT NULL,
    studio_id UUID NOT NULL,
    seat_type_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_seats_studios FOREIGN KEY (studio_id) REFERENCES studios(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_seats_seat_type FOREIGN KEY (seat_type_id) REFERENCES seat_type(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);