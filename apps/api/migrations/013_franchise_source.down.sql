DROP INDEX IF EXISTS idx_franchises_category;
DROP INDEX IF EXISTS idx_franchises_source_external;
ALTER TABLE franchises
    DROP COLUMN IF EXISTS category,
    DROP COLUMN IF EXISTS external_id,
    DROP COLUMN IF EXISTS source;
