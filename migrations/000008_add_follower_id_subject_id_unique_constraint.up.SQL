ALTER TABLE public.follow
ADD CONSTRAINT unique_follower_subject_pair UNIQUE(follower_id, subject_id),
ADD CONSTRAINT check_different_ids CHECK(follower_id <> subject_id);