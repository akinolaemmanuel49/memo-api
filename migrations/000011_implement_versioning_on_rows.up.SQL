ALTER TABLE public.users ADD COLUMN _version INTEGER DEFAULT 0;
ALTER TABLE public.memos ADD COLUMN _version INTEGER DEFAULT 0;
ALTER TABLE public.follow ADD COLUMN _version INTEGER DEFAULT 0;
ALTER TABLE public.likes ADD COLUMN _version INTEGER DEFAULT 0;