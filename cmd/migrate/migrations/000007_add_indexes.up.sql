CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING GIN (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING GIN (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING GIN (tags);


-- index on users table

CREATE INDEX IF NOT EXISTS idx_users_name ON users (username);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);