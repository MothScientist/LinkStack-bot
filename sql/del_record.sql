UPDATE "links"
SET "Status" = FALSE
WHERE
    "TelegramId" = ?
    AND "LinkId" = ?
    AND "Status" IS TRUE