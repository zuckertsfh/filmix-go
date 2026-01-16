CREATE TABLE transactions (
    id UUID NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL,
    external_ref VARCHAR(255) NOT NULL UNIQUE,
    invoice_number TEXT,
    amount BIGINT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ,
    payment_method_id UUID NOT NULL,
    showtime_id UUID NOT NULL,
    theater_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_transactions_payment_method FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_transactions_showtime FOREIGN KEY (showtime_id) REFERENCES showtimes(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_transactions_theater FOREIGN KEY (theater_id) REFERENCES theaters(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);