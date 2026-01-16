CREATE TABLE studios (
    id UUID NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    theater_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_studios_theaters FOREIGN KEY (theater_id) REFERENCES theaters(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);