CREATE TABLE films (
    id int NOT NULL AUTO_INCREMENT,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title varchar(100) NOT NULL,
    director varchar(100) NOT NULL,
    release_date date NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT UC_FILM_TITLE UNIQUE (title)
);