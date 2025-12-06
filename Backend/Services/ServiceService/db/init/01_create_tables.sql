CREATE TABLE user_service_profiles (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    service         TEXT   NOT NULL,
    provider_user_id TEXT  NOT NULL,
    access_token    TEXT   NOT NULL,
    refresh_token   TEXT,
    expires_at      TIMESTAMPTZ,
    raw_profile     JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, service),
    UNIQUE (service, provider_user_id)
);

CREATE INDEX idx_user_service_profiles_user_id
    ON user_service_profiles(user_id);

CREATE INDEX idx_user_service_profiles_user_service
    ON user_service_profiles(user_id, service);

CREATE TABLE user_service_fields IF NOT EXISTS (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    profile_id      BIGINT NOT NULL REFERENCES user_service_profiles(id) ON DELETE CASCADE,
    field_key       TEXT   NOT NULL,
    value_string    TEXT,
    value_number    DOUBLE PRECISION,
    value_boolean   BOOLEAN,
    value_json      JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (profile_id, field_key)
);

CREATE INDEX idx_user_service_fields_profile
    ON user_service_fields(profile_id);

CREATE INDEX idx_user_service_fields_service_key_value
    ON user_service_fields(field_key, value_string, profile_id);