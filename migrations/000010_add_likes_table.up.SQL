-- noinspection SpellCheckingInspectionForFile

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- noinspection SqlResolve
CREATE TABLE public.likes
(
    id         UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    liked_by   UUID        NOT NULL,
    memo_id    UUID        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (liked_by) REFERENCES public.users (id),
    FOREIGN KEY (memo_id) REFERENCES public.memos (id),
    CONSTRAINT unique_liked_by_memo_id_pair UNIQUE (liked_by, memo_id)
);
