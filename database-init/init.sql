CREATE TABLE links (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    TelegramId INTEGER NOT NULL,
    LinkId INTEGER NOT NULL,
    Status INTEGER NOT NULL,
    Url TEXT NOT NULL,
    Title TEXT NOT NULL,
    Date TEXT DEFAULT (CURRENT_TIMESTAMP)
);

-- Index to search all user records
CREATE INDEX idx_tg_id ON links(TelegramId);

-- Index for searching valid user records
CREATE INDEX idx_tg_id_active_links ON links(TelegramId) WHERE Status = 1;

-- Composite index for searching a specific user record
CREATE INDEX idx_tg_link ON links(TelegramId, LinkId);
