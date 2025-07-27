INSERT INTO "links" (
    "TelegramId",
    "LinkId",
    "Status",
    "Url",
    "Title"
)
SELECT
    ?,
    (SELECT COALESCE(MAX("LinkId"), 0) + 1 FROM "links" WHERE "TelegramId" = ?),
    TRUE,
    ?,
    ?
RETURNING "LinkId";