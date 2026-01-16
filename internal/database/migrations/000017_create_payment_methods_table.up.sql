CREATE TABLE payment_methods (
    id UUID NOT NULL UNIQUE,
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    logo_url TEXT NOT NULL,
    active BOOLEAN NOT NULL,
    payment_method_type_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_payment_methods_type FOREIGN KEY (payment_method_type_id) REFERENCES payment_method_types(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);