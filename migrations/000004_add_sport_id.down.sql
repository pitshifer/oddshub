ALTER TABLE events ADD COLUMN sport TEXT;

UPDATE events
SET sport = sports.key
FROM sports
WHERE sports.id = events.sport_id;

ALTER TABLE events ALTER COLUMN sport SET NOT NULL;

DROP INDEX IF EXISTS idx_events_sport_id;
ALTER TABLE events DROP COLUMN sport_id;
CREATE INDEX idx_events_sport ON events(sport);
