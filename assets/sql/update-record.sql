UPDATE records
	SET
    start = $1,
    "end" = $2,
    actualStart = $3,
    actualEnd = $4,
		"group" = $5,
    allDay = $6,
    title = $7,
    url = $8,
    description = $9,
    pid = $10,
		priority = $11,
    metadata = $12,
    status = $13
WHERE id = $14;
