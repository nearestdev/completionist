ALTER TABLE franchises
    ADD COLUMN IF NOT EXISTS source VARCHAR(50),
    ADD COLUMN IF NOT EXISTS external_id VARCHAR(255),
    ADD COLUMN IF NOT EXISTS category VARCHAR(50);

CREATE UNIQUE INDEX IF NOT EXISTS idx_franchises_source_external ON franchises(source, external_id)
    WHERE source IS NOT NULL AND external_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_franchises_category ON franchises(category)
    WHERE category IS NOT NULL;
