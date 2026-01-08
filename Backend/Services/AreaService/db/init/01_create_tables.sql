CREATE TABLE IF NOT EXISTS areas (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
    active BOOLEAN NOT NULL,
    user_id INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS areas_user_id_idx ON areas (user_id);

CREATE TABLE IF NOT EXISTS reactions (
    id SERIAL PRIMARY KEY,
    area_id SERIAL NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    service TEXT NOT NULL,
    title TEXT NOT NULL,
    inputs JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS reactions_area_id_idx ON reactions (area_id);

CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    area_id SERIAL NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    service TEXT NOT NULL,
    title TEXT NOT NULL,
    inputs JSONB NOT NULL,
    type TEXT NOT NULL
)