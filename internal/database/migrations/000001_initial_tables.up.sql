BEGIN;

CREATE TABLE IF NOT EXISTS artists (
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

CREATE TABLE IF NOT EXISTS locations (
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

CREATE TABLE IF NOT EXISTS events (
                                     id              SERIAL PRIMARY KEY,
                                     name            TEXT,
                                    start_time       TIMESTAMPTZ,
                                     location        SERIAL REFERENCES locations
);

CREATE TABLE IF NOT EXISTS invited_artists (
                                             artist_id       SERIAL REFERENCES artists,
                                             event_id        SERIAL REFERENCES events,
                                             travel_expenses MONEY,
                                             confirmation    TEXT
);

CREATE TABLE IF NOT EXISTS artworks (
                                       id                  SERIAL PRIMARY KEY,
                                       title               TEXT,
                                       artist_id           SERIAL REFERENCES artists,
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

CREATE TABLE IF NOT EXISTS artwork_event_locations (
                                                    artwork_id                      SERIAL REFERENCES artworks,
                                                    event_id                        SERIAL REFERENCES events,
                                                    location_id                     SERIAL REFERENCES locations,
                                                    will_be_sent_by_post            BOOLEAN,
                                                    will_be_sent_by_spedition       BOOLEAN,
                                                    is_collected_after_exhibition   BOOLEAN,                                -- whether artist picks up items themself after the fact
                                                    is_built_onsite                 BOOLEAN,                                -- whether item comes prebuilt or not
                                                    is_built_by_artist              BOOLEAN,
                                                    shipping_address_id             SERIAL REFERENCES locations,             -- place where the artwork is shipped from to the event
                                                    packaging                       TEXT,
                                                    material                        TEXT,
                                                    no_pieces                       INTEGER,
                                                    size                            FLOAT,
                                                    pub_agreement                   TEXT
);

COMMIT;