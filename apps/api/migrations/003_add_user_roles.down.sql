DROP INDEX IF EXISTS idx_users_role;

ALTER TABLE users
    DROP COLUMN IF EXISTS role;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        DROP TYPE user_role;
    END IF;
END
$$;
