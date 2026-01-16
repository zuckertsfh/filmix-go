CREATE TABLE movie_statuses (
    id UUID NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL,
    PRIMARY KEY(id)
);