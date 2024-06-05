-- noinspection SpellCheckingInspectionForFile

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- noinspection SqlResolve
CREATE TABLE public.users (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username        VARCHAR(255) NOT NULL UNIQUE,
    first_name      VARCHAR(255) NOT NULL,
    last_name       VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    avatar          TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT '',
    about           TEXT NOT NULL DEFAULT '',
    is_activated    BOOLEAN NOT NULL DEFAULT FALSE,
    follower_count  INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted         BOOLEAN NOT NULL DEFAULT FALSE,
    _version        INTEGER DEFAULT 0,

    CONSTRAINT users_username_length_check CHECK (length(username) >= 5),
    CONSTRAINT users_first_name_length_check CHECK (length(first_name) >= 2 AND length(first_name) <= 100),
    CONSTRAINT users_last_name_length_check CHECK (length(last_name) >= 2 AND length(last_name) <= 100)
);