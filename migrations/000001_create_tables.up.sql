CREATE TABLE providers (
    id   SMALLSERIAL PRIMARY KEY,
    key  TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE events (
    id          TEXT        NOT NULL,
    provider_id SMALLINT    NOT NULL REFERENCES providers(id),
    sport       TEXT        NOT NULL,
    home_team   TEXT        NOT NULL,
    away_team   TEXT        NOT NULL,
    start_time  TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, provider_id)
);

CREATE TABLE bookmakers (
    id  SERIAL PRIMARY KEY,
    key TEXT UNIQUE NOT NULL
);

CREATE TABLE odds (
    id           BIGSERIAL   PRIMARY KEY,
    event_id     TEXT        NOT NULL,
    provider_id  SMALLINT    NOT NULL,
    bookmaker_id INT         NOT NULL REFERENCES bookmakers(id),
    market       TEXT        NOT NULL,
    outcome      TEXT        NOT NULL,
    price        NUMERIC(10, 4) NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (event_id, provider_id) REFERENCES events(id, provider_id)
);

-- Быстрый поиск котировок по событию
CREATE INDEX idx_odds_event ON odds(event_id, provider_id);

-- Для будущего партиционирования и очистки старых данных
CREATE INDEX idx_odds_collected_at ON odds(collected_at);

-- Фильтрация событий по виду спорта
CREATE INDEX idx_events_sport ON events(sport);
