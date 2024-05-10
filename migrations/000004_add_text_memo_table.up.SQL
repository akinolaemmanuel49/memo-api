-- noinspection SpellCheckingInspectionForFile

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- noinspection SqlResolve
CREATE TABLE public.text_memos (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    content         VARCHAR(255) NOT NULL,
    likes           INTEGER NOT NULL DEFAULT 0,
    shares          INTEGER NOT NULL DEFAULT 0,
    owner_id        UUID,
    deleted         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY(owner_id) REFERENCES public.users(id)
);