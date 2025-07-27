SELECT COALESCE("LinkId", 0)
FROM "links"
WHERE
    "TelegramId" = ?
    AND "Url" = ?
    AND "Status" IS TRUE