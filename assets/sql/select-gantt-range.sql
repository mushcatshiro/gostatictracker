SELECT
	MIN("start") as startDate,
	MAX("end") as endDate
FROM records
WHERE "group" = $1
