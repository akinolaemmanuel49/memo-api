-- noinspection SqlResolve
CREATE TABLE public.memos (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id        UUID NOT NULL,
    type            VARCHAR(50) NOT NULL,
    content         TEXT,
    likes           INTEGER NOT NULL DEFAULT 0,
    shares          INTEGER NOT NULL DEFAULT 0,
    resource_url    TEXT,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted         BOOLEAN NOT NULL DEFAULT FALSE,
    _version        INTEGER DEFAULT 0,


    FOREIGN KEY(owner_id) REFERENCES public.users(id)
);
