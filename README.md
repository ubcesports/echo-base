# echo-base
echo-base serves as the backend for all services for UBC Esports. 

## Running the application
To run the application, create a `.env` file following `.env.example`.

Then run:
```
docker compose up --build -d
```

## Setting up for development
The application can be run in development using either docker or manually.
The docker setup is preferred because it does not require additional
configuration and will automatically spin up a database for you.

### Development Environment 
echo-base uses Docker for its development environment. 
This is the recommended way to setup the application.
To run the application using Docker, run 
```
docker compose -f compose.dev.yaml up -d
``` 
The application will be accessible at `http://localhost:8080` and the database can be found at `http://localhost:5432`.

You can log into the database using these credentials: `postgresql://user:pass@localhost/echobase?sslmode=disable`.

echo-base uses `sql-migrate` to manage its database migrations. You can read the `sql-migrate` docs for more information. 
echo-base also uses Go's tool dependencies for managing development tools. 
You can apply the migrations by running
```
go tool sql-migrate up
```
If there is a need to reset the database, run `docker volume rm echo-base_postgres-data`. 

#### Docker Setup with Live Reload
The development environment docker compose uses [air](https://github.com/air-verse/air) to automatically watch for
source file changes and reload the application.

It is important to note that the development environment uses watch mode instead
of the normal docker compose file found at `compose.yaml`. When using watch mode 
during development, make sure to test the application with the full docker build script prior to merging with `main`,
because the regular build script has various build and container optimizations 
that are not used in watch mode.

### Manual Setup
If you wish to setup your development environment manually:
Create a `.env` file in `cmd/server`. Ask your directors for the specific
values for the database connection. 

To start the app, `cd` into `cmd/server` and run `go run .`. 

To check that the server is working, you can run a health ping through the
`GET /health` endpoint. 

### Testing
