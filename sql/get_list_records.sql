SELECT
    "LinkId",
    "Url",
    "Title"
FROM "links"
WHERE
    "TelegramId" = ?
    AND "Status" IS TRUE