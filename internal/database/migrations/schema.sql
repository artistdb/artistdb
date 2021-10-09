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

CREATE TABLE IF NOT EXISTS Location (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR NOT NULL,
    country         VARCHAR,
    zip             VARCHAR,
    city            VARCHAR,
    street          VARCHAR,
    picture         VARCHAR,
    description     VARCHAR,
    lat             VARCHAR,
    lon             VARCHAR
);

CREATE TABLE IF NOT EXISTS Event (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR,
    date            TIMESTAMP,
    location        SERIAL REFERENCES Location
);

CREATE TYPE confirmation AS ENUM ('not confirmed', 'confirmed');

CREATE TABLE IF NOT EXISTS InvitedArtist (
    artist_id       SERIAL REFERENCES Artist,
    event_id        SERIAL REFERENCES Event,
    travel_expenses MONEY,
    confirmation    confirmation
);

CREATE TABLE IF NOT EXISTS Artwork (
    id                  SERIAL PRIMARY KEY,
    title               VARCHAR,
    artist_id           SERIAL REFERENCES Artist,
    synopsis_en         VARCHAR,
    synopsis_ger        VARCHAR,
    picture_1           VARCHAR,
    picture_2           VARCHAR,
    picture_3           VARCHAR,
    material_demands    VARCHAR,
    insurance_amount    MONEY,
    sales_val           MONEY,
    height              FLOAT,
    length              FLOAT,
    width               FLOAT,
    weight              FLOAT,
    category            VARCHAR
);

CREATE TABLE IF NOT EXISTS ArtworkEventLocation (
    artwork_id          SERIAL REFERENCES Artwork,
    event_id            SERIAL REFERENCES Event,
    location_id         SERIAL REFERENCES Location,
    by_post             BOOLEAN,
    by_spedition        BOOLEAN,
    is_collected        BOOLEAN,
    is_built_onsite     BOOLEAN,
    is_built_by_artist  BOOLEAN,
    address_id          SERIAL REFERENCES Location,
    packaging           VARCHAR,
    material            VARCHAR,
    no_pieces           INTEGER,
    size                FLOAT,
    pub_agreement       VARCHAR 
);