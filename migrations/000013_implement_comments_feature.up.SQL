-- noinspection SpellCheckingInspectionForFile

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- noinspection SqlResolve
CREATE TABLE public.comments
(
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id UUID REFERENCES public.users (id),
    memo_id UUID REFERENCES public.memos (id),
    parent_id UUID REFERENCES public.comments (id),
    comment_type VARCHAR(50) NOT NULL,
    comment_content TEXT,
    caption TEXT,
    transcript TEXT,
    likes INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted BOOLEAN DEFAULT FALSE,
    _version INTEGER DEFAULT 0
);