CREATE TABLE IF NOT EXISTS records (
  "id" INTEGER PRIMARY KEY,
  start TIMESTAMP,
  "end" TIMESTAMP,
  actualStart TIMESTAMP,
  actualEnd TIMESTAMP,
  insertTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "group" TEXT DEFAULT 'ungroupped' NOT NULL,
  defaultMode TEXT DEFAULT '',
  repeat TEXT DEFAULT '',
  allDay BOOLEAN DEFAULT FALSE,
  title TEXT NOT NULL,
  url TEXT,
  description TEXT,
  pid INTEGER,
  priority INT,
  metadata TEXT,
  status INT
)
