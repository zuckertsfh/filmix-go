CREATE TABLE theaters (
    id UUID NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(9,6) NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    cinema_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_theaters_cinemas FOREIGN KEY (cinema_id) REFERENCES cinemas(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);