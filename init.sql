CREATE TABLE IF NOT EXISTS data (
                                    short_url TEXT PRIMARY KEY,
                                    long_url TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS short_url_hash_idx ON data USING hash (short_url);