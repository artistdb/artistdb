BEGIN;

CREATE TABLE IF NOT EXISTS artists (
                                      id              UUID PRIMARY KEY,
                                      first_name      TEXT NOT NULL,
                                      last_name       TEXT NOT NULL,
                                      artist_name     TEXT,
                                      pronouns        TEXT[],
                                      date_of_birth   TIMESTAMPTZ,
                                      place_of_birth  TEXT,
                                      nationality     TEXT,
                                      language        TEXT,
                                      facebook        TEXT,
                                      instagram       TEXT,
                                      bandcamp        TEXT,
                                      bio_ger         TEXT,
                                      bio_en          TEXT,
                                      created_at      TIMESTAMPTZ NOT NULL,
                                      updated_at      TIMESTAMPTZ NOT NULL,
                                      deleted_at      TIMESTAMPTZ,
                                      email           TEXT
);

CREATE TABLE IF NOT EXISTS locations (
                                        id              UUID PRIMARY KEY,
                                        created_at      TIMESTAMPTZ NOT NULL,
                                        updated_at      TIMESTAMPTZ NOT NULL,
                                        deleted_at      TIMESTAMPTZ,
                                        name            TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
                                     id              UUID PRIMARY KEY,
                                     created_at      TIMESTAMPTZ NOT NULL,
                                     updated_at      TIMESTAMPTZ NOT NULL,
                                     deleted_at      TIMESTAMPTZ,
                                     name            TEXT NOT NULL,
                                     start_time      TIMESTAMPTZ,
                                     location_id    UUID,
                                     CONSTRAINT     fk_locations
                                        FOREIGN KEY (location_id)
                                            REFERENCES locations(id)
                                            ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS artist_event (
                                            artist_id   UUID REFERENCES artists ON UPDATE CASCADE ON DELETE CASCADE,
                                            event_id    UUID REFERENCES events ON UPDATE CASCADE ON DELETE CASCADE,
                                            confirmed   BOOL DEFAULT false,
                                            CONSTRAINT  artist_event_pk PRIMARY KEY (artist_id, event_id)
);

-- CREATE TABLE IF NOT EXISTS artworks (
--                                        id                  UUID PRIMARY KEY,
--                                        title               TEXT,
--                                        artist_id           UUID REFERENCES artists,
--                                        synopsis_en         TEXT,
--                                        synopsis_ger        TEXT,
--                                        pictures            TEXT[],
--                                        material_demands    TEXT,
--                                        insurance_amount    MONEY,
--                                        sales_val           MONEY,
--                                        height              FLOAT,
--                                        length              FLOAT,
--                                        width               FLOAT,
--                                        weight              FLOAT,
--                                        category            TEXT
-- );
--
-- CREATE TABLE IF NOT EXISTS artwork_event_locations (
--                                                     artwork_id                      UUID REFERENCES artworks,
--                                                     event_id                        UUID REFERENCES events,
--                                                     location_id                     UUID REFERENCES locations,
--                                                     will_be_sent_by_post            BOOLEAN,
--                                                     will_be_sent_by_spedition       BOOLEAN,
--                                                     is_collected_after_exhibition   BOOLEAN,                                -- whether artist picks up items themself after the fact
--                                                     is_built_onsite                 BOOLEAN,                                -- whether item comes prebuilt or not
--                                                     is_built_by_artist              BOOLEAN,
--                                                     shipping_address_id             UUID REFERENCES locations,             -- place where the artwork is shipped from to the event
--                                                     packaging                       TEXT,
--                                                     material                        TEXT,
--                                                     no_pieces                       INTEGER,
--                                                     size                            FLOAT,
--                                                     pub_agreement                   TEXT
-- );

COMMIT;