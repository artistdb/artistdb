BEGIN;

DROP TABLE IF EXISTS artists CASCADE;
DROP TABLE IF EXISTS locations CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS invited_artists CASCADE;
DROP TABLE IF EXISTS artworks CASCADE;
DROP TABLE IF EXISTS artwork_event_locations CASCADE;

COMMIT;