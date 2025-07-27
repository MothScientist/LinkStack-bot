SELECT
    COALESCE("Url", ''),
    COALESCE("Title", ''),
    COALESCE("Status", FALSE)
FROM "links"
WHERE
    "TelegramId" = ?
    AND "LinkId" = ?