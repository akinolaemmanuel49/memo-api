-- noinspection SpellCheckingInspectionForFile

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- noinspection SqlResolve
CREATE TABLE public.shares
(
    id         UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    shared_by  UUID        NOT NULL,
    memo_id    UUID        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    _version   INTEGER              DEFAULT 0,
    FOREIGN KEY (shared_by) REFERENCES public.users (id),
    FOREIGN KEY (memo_id) REFERENCES public.memos (id),
    CONSTRAINT unique_shared_by_memo_id_pair UNIQUE (shared_by, memo_id)
);
