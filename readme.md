# GoStaticTracker

after sometime into development, the author realized that this project is just
a spiritual successor of a previous project's submodule `Negotium`. Intends to
support 4 type of rendering, gantt, todo list, calendar and overall statistics.

## Start

### Getting PostgreSQL Up

```bash
>>> docker volume create pgsql-data  # recommended otherwise no persistent data
>>> docker run -d \
  --name pgsql-toy-db \
  -p 5432:5432 \
  -e POSTGRES_USER=pgsql \
  -e POSTGRES_PASSWORD=pgsql \
  -e POSTGRES_DB=pgsql \
  -v pgsql-data:/var/lib/postgresql/data \
  postgres:16-alpine
```

## Run Test

```bash
>>> cd ~/gostatictracker
>>> docker run -d --rm --name pgsql-toy-db -p 5432:5432 psql-toy-db  # create DB
>>> go run mock/main.go pgsql pgsql localhost pgsql  # create DB entries
>>> go test ./render -v

>>> go run cli/main.go --outputDir "E:\\statictracker\\output"  # windows path
```
