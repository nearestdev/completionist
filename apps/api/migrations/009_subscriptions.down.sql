DROP INDEX IF EXISTS idx_subscriptions_current_period_end;
DROP INDEX IF EXISTS idx_subscriptions_stripe_subscription_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP TABLE IF EXISTS subscriptions;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'subscription_status') THEN
        DROP TYPE subscription_status;
    END IF;
END
$$;

-- Note: Cannot remove 'member' from user_role ENUM in PostgreSQL.
-- The value will remain but be unused after rollback.
