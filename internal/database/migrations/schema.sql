CREATE DATABASE main;

CREATE TABLE IF NOT EXISTS Artist (
    id              SERIAL PRIMARY KEY,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    artist_name     TEXT,
    pronouns        TEXT,
    date_of_birth   TIMESTAMPTZ,
    place_of_birth  TEXT,
    nationality     TEXT,
    language        TEXT,
    facebook        TEXT,
    instagram       TEXT,
    bandcamp        TEXT,
    bio_ger         TEXT,
    bio_en          TEXT
);

CREATE TABLE IF NOT EXISTS Location (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    country         TEXT,
    zip             TEXT,
    city            TEXT,
    street          TEXT,
    picture         TEXT,
    description     TEXT,
    lat             TEXT,
    lon             TEXT
);

CREATE TABLE IF NOT EXISTS Event (
    id              SERIAL PRIMARY KEY,
    name            TEXT,
    date            TIMESTAMPTZ,
    location        SERIAL REFERENCES Location
);

CREATE TABLE IF NOT EXISTS InvitedArtist (
    artist_id       SERIAL REFERENCES Artist,
    event_id        SERIAL REFERENCES Event,
    travel_expenses MONEY,
    confirmation    TEXT
);

CREATE TABLE IF NOT EXISTS Artwork (
    id                  SERIAL PRIMARY KEY,
    title               TEXT,
    artist_id           SERIAL REFERENCES Artist,
    synopsis_en         TEXT,
    synopsis_ger        TEXT,
    picture_1           TEXT,
    picture_2           TEXT,
    picture_3           TEXT,
    material_demands    TEXT,
    insurance_amount    MONEY,
    sales_val           MONEY,
    height              FLOAT,
    length              FLOAT,
    width               FLOAT,
    weight              FLOAT,
    category            TEXT
);

CREATE TABLE IF NOT EXISTS ArtworkEventLocation (
    artwork_id                      SERIAL REFERENCES Artwork,
    event_id                        SERIAL REFERENCES Event,
    location_id                     SERIAL REFERENCES Location,
    will_be_sent_by_post            BOOLEAN,
    will_be_sent_by_spedition       BOOLEAN,
    is_collected_after_exhibition   BOOLEAN,                                -- whether artist picks up items themself after the fact
    is_built_onsite                 BOOLEAN,                                -- whether item comes prebuilt or not
    is_built_by_artist              BOOLEAN,
    shipping_address_id             SERIAL REFERENCES Location,             -- place where the artwork is shipped from to the event
    packaging                       TEXT,
    material                        TEXT,
    no_pieces                       INTEGER,
    size                            FLOAT,
    pub_agreement                   TEXT 
);