-- noinspection SqlResolve
CREATE TABLE public.memos (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id        UUID NOT NULL,
    memo_type       VARCHAR(50) NOT NULL,
    memo_content    TEXT,
    likes           INTEGER NOT NULL DEFAULT 0,
    shares          INTEGER NOT NULL DEFAULT 0,
    caption         TEXT,
    transcript      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted         BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY(owner_id) REFERENCES public.users(id)
);

-- Copy data from text_memos to memos
-- noinspection SqlResolve
INSERT INTO public.memos (id, owner_id, memo_type, memo_content, likes, shares, caption, transcript, created_at, updated_at)
SELECT id, owner_id, 'text', content, likes, shares, '', '', created_at, updated_at
FROM public.text_memos;

-- Drop the text_memos table
-- noinspection SqlResolve
DROP TABLE public.text_memos;
