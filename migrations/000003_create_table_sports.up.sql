CREATE TABLE sports (
    id   SMALLSERIAL PRIMARY KEY,
    key  VARCHAR(250) UNIQUE NOT NULL,
    title VARCHAR(250) NOT NULL,
    group_name VARCHAR(250),
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    has_outrights BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_sports_key ON sports(key);

CREATE INDEX idx_sports_group ON sports(group_name);
