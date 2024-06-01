-- noinspection SqlResolve
CREATE TABLE public.text_memos (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    content         VARCHAR(255) NOT NULL,
    likes           INTEGER NOT NULL DEFAULT 0,
    shares          INTEGER NOT NULL DEFAULT 0,
    owner_id        UUID NOT NULL,
    deleted         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY(owner_id) REFERENCES public.users(id)
);

-- Copy data from memos back to text_memos
INSERT INTO public.text_memos (id, content, likes, shares, owner_id, created_at, updated_at)
SELECT id, memo_content, likes, shares, owner_id, created_at, updated_at
FROM public.memos WHERE memo_type = 'text';

-- Drop the memos table
DROP TABLE public.memos;
