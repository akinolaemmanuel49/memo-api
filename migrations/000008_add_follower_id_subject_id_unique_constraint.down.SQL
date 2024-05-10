ALTER TABLE public.follow
DROP CONSTRAINT unique_follower_subject_pair,
DROP CONSTRAINT check_different_ids;