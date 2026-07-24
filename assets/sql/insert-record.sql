INSERT INTO records (
  start,
  "end",
  actualStart,
  actualEnd,
  "group",
  defaultMode,
  repeat,
  allDay,
  title,
  url,
  description,
  pid,
  priority,
  metadata,
  status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id;
