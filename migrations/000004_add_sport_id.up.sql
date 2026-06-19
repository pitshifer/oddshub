-- Add sport_id column to events table
ALTER TABLE events
    ADD COLUMN sport_id SMALLINT REFERENCES sports(id);

-- Update existing records to set sport_id based on the sport column
UPDATE events
SET sport_id = sports.id
FROM sports WHERE sports.key = events.sport;

-- ALTER TABLE events
--     ALTER COLUMN sport_id SET NOT NULL;

DROP INDEX IF EXISTS idx_events_sport;

ALTER TABLE events DROP COLUMN sport;

CREATE INDEX idx_events_sport_id ON events(sport_id);
