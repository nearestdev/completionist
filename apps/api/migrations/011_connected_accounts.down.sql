DROP INDEX IF EXISTS idx_import_jobs_created_at;
DROP INDEX IF EXISTS idx_import_jobs_connected_account;
DROP INDEX IF EXISTS idx_import_jobs_status;
DROP INDEX IF EXISTS idx_import_jobs_user_id;
DROP TABLE IF EXISTS import_jobs;

DROP INDEX IF EXISTS idx_connected_accounts_provider_user;
DROP INDEX IF EXISTS idx_connected_accounts_provider;
DROP INDEX IF EXISTS idx_connected_accounts_user_id;
DROP TABLE IF EXISTS connected_accounts;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'import_job_status') THEN
        DROP TYPE import_job_status;
    END IF;
END
$$;
