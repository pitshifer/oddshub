CREATE TABLE odds_lines (
    id           BIGSERIAL PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    bookmaker_id INT NOT NULL REFERENCES bookmakers(id) ON DELETE RESTRICT,
    market       VARCHAR(50) NOT NULL,
    outcome      VARCHAR(100) NOT NULL,
    price        NUMERIC(10, 4) NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (event_id, bookmaker_id, market, outcome)
);

CREATE TABLE odds_price_history (
    id      BIGSERIAL PRIMARY KEY,
    line_id BIGINT NOT NULL REFERENCES odds_lines(id) ON DELETE CASCADE,
    price   NUMERIC(10, 4) NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_odds_price_history_line ON odds_price_history(line_id, collected_at);
