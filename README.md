# echo-base

## Setting up for development
The application can be run in development using either docker or manually.
The docker setup is preferred because it does not require additional
configuration and will automatically spin up a database for you.

### Docker Setup (recommended)
To run the application using docker, run the command 
`docker compose up --build` in a new shell. The application should now be 
accessible at `http://localhost:8080`, and the database should be accessible
at `postgresql://user:pass@localhost/echobase?sslmode=disable`.

Once up, ensure your local database has up-to-date migrations by running
`go tool sql-migrate up`. If there is a need to reset the database,
run `docker volume rm echo-base_postgres-data` while the application is
not running.

### Docker Setup with Live Reload
During development, it is possible to use the development compose configuration
which uses [air](https://github.com/air-verse/air) to automatically watch for
source file changes and reload the application.

To use watch mode, run `docker compose -f compose.dev.yaml up` instead of the
normal up command. When using watch mode during development, make sure to test
the application with the full docker build script prior to merging with `main`,
because the regular build script has various build and container optimizations 
that are not used in watch mode.

### Manual Setup
Create a `.env` file in `cmd/server`. Ask your directors for the specific
values for the database connection. 

To start the app, `cd` into `cmd/server` and run `go run .`. 

To check that the server is working, you can run a health ping through the
`GET /health` endpoint. 
