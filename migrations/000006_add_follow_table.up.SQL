-- noinspection SpellCheckingInspectionForFile

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- noinspection SqlResolve
CREATE TABLE public.follow (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    follower_id     UUID NOT NULL,
    subject_id      UUID NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY(follower_id) REFERENCES public.users(id),
    FOREIGN KEY(subject_id) REFERENCES public.users(id)
)