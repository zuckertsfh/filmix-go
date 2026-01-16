CREATE TABLE payment_method_types (
    id UUID NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY(id)
);