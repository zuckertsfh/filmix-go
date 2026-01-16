CREATE TABLE transaction_items (
    id UUID NOT NULL UNIQUE,
    price BIGINT NOT NULL,
    transaction_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    seat_type_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_items_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_items_seat FOREIGN KEY (seat_id) REFERENCES seats(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_items_seat_type FOREIGN KEY (seat_type_id) REFERENCES seat_type(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);