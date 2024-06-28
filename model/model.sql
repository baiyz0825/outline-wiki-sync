CREATE TABLE IF NOT EXISTS file_sync_record
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    outline_wiki_id TEXT,
    collection_id   TEXT,
    file_name       TEXT,
    file_size       TEXT,
    file_path       TEXT,
    file_content    TEXT,
    sync            INTEGER,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted         INTEGER
);

CREATE INDEX IF NOT EXISTS idx_outline_wiki_id ON file_sync_record(outline_wiki_id);
CREATE INDEX IF NOT EXISTS idx_collection_id ON file_sync_record(collection_id);
CREATE INDEX IF NOT EXISTS idx_file_name ON file_sync_record(file_name);

CREATE TABLE IF NOT EXISTS outline_wiki_collection_mapping
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_id   TEXT,
    current_id      TEXT,
    parent_id       TEXT,
    collection_path TEXT,
    collection_name TEXT,
    sync            INTEGER,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted         INTEGER
);

CREATE INDEX IF NOT EXISTS idx_collection_id ON outline_wiki_collection_mapping(collection_id);
CREATE INDEX IF NOT EXISTS idx_collection_path ON outline_wiki_collection_mapping(collection_path);
CREATE INDEX IF NOT EXISTS idx_collection_name ON outline_wiki_collection_mapping(collection_name);