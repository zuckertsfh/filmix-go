CREATE TABLE seat_type (
    id UUID NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    cinema_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_seat_type_cinema FOREIGN KEY (cinema_id) REFERENCES cinemas(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);