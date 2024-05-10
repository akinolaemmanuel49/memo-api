ALTER TABLE public.users
ALTER COLUMN username TYPE VARCHAR(255),
ALTER COLUMN username SET NOT NULL,
ADD CONSTRAINT users_username_key UNIQUE (username);