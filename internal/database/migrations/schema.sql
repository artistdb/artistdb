CREATE DATABASE main;

CREATE TYPE pronouns AS ENUM ('they', 'he', 'she');

CREATE TABLE IF NOT EXISTS Artist (
    id              SERIAL PRIMARY KEY,
    first_name      VARCHAR NOT NULL,
    last_name       VARCHAR NOT NULL,
    artist_name     VARCHAR,
    pronouns        pronouns,
    date_of_birth   TIMESTAMP,
    place_of_birth  VARCHAR,
    nationality     VARCHAR,
    language        VARCHAR,
    facebook        VARCHAR,
    instagram       VARCHAR,
    bandcamp        VARCHAR,
    bio_ger         VARCHAR,
    bio_en          VARCHAR
);