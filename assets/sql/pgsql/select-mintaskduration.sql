SELECT MIN(("end"::date - start::date)+1) AS minTaskDuration FROM events WHERE "group" = $1
