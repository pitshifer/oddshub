DROP TABLE IF EXISTS odds;
DROP TABLE IF EXISTS events;

CREATE TABLE events (
    id          TEXT        NOT NULL,
    provider_id SMALLINT    NOT NULL REFERENCES providers(id),
    sport       TEXT        NOT NULL,
    home_team   TEXT        NOT NULL,
    away_team   TEXT        NOT NULL,
    start_time  TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, provider_id)
);

CREATE TABLE odds (
    id           BIGSERIAL      PRIMARY KEY,
    event_id     TEXT           NOT NULL,
    provider_id  SMALLINT       NOT NULL,
    bookmaker_id INT            NOT NULL REFERENCES bookmakers(id),
    market       TEXT           NOT NULL,
    outcome      TEXT           NOT NULL,
    price        NUMERIC(10, 4) NOT NULL,
    collected_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    FOREIGN KEY (event_id, provider_id) REFERENCES events(id, provider_id)
);

CREATE INDEX idx_odds_event      ON odds(event_id, provider_id);
CREATE INDEX idx_odds_collected_at ON odds(collected_at);
CREATE INDEX idx_events_sport    ON events(sport);
