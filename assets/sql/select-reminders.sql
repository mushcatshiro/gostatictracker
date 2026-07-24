WITH month_range AS (
	SELECT tsrange(
		make_timestamp($1, $2, 1, 0, 0, 0),
		make_timestamp($1, $2, 1, 0, 0, 0) + interval '1 month',
		'[)'
	) AS period
)
SELECT
  title,
  "group",
  CASE
    WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN "actualstart"
    WHEN "start" IS NOT NULL OR "end" IS NOT NULL THEN "start"
    ELSE "inserttime"
  END as s,
  CASE
    WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN "actualend"
    WHEN "start" IS NOT NULL OR "end" IS NOT NULL THEN "end"
    ELSE "inserttime"
  END as e
FROM events, month_range
WHERE
  (CASE
    WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN
      tsrange("actualstart", "actualend", '[]') && month_range.period
    WHEN "start" IS NOT NULL OR "end" IS NOT NULL THEN
      tsrange("start", "end", '[]') && month_range.period
    ELSE
      "inserttime" >= LOWER(month_range.period) AND "inserttime" < upper(month_range.period)
  END)
  AND
    (CASE
      WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN
        ("actualend" - "actualstart") <= interval '7 days'
      WHEN ("actualstart" IS NULL OR "actualend" IS NULL) AND ("start" IS NOT NULL AND "end" IS NOT NULL) THEN
        ("end" - "start") <= interval '7 days'
      ELSE
        TRUE
    END);
