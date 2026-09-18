# Run URL App

URL App is a Go API, PostgreSQL database, and React/Vite frontend. The easiest way to run the complete application locally is with Podman Compose (Docker Compose commands are equivalent if Docker is installed).

## Prerequisites

- Git
- Podman and `podman-compose` (or Docker and Compose)
- For native development: Go 1.27+, Node.js/npm, and PostgreSQL

## Run with containers

From the repository root:

```sh
cp .env.example .env
```

Edit `.env` and set a password. For local HTTP, use:

```env
POSTGRES_PASSWORD=choose-a-long-random-password
VITE_SHORT_URL_BASE=http://localhost:8080
APP_PORT=8080
```

Build and start the database, API, and web UI:

```sh
podman compose up -d --build
podman compose ps
```

Open <http://localhost:8080>. Check the API health endpoint with:

```sh
curl http://127.0.0.1:8080/health
```

View logs when diagnosing startup or request errors:

```sh
podman compose logs -f --tail=100
```

Stop the services while keeping database data:

```sh
podman compose down
```

The database schema is mounted from `backend/migrations/001_create_urls.sql` and is applied automatically when the PostgreSQL volume is first created. To intentionally remove all database data, use `podman compose down -v`.

## Run the production HTTPS stack

Set these values in `.env`:

```env
POSTGRES_PASSWORD=<long-random-password>
DOMAIN=go.example.com
VITE_SHORT_URL_BASE=https://go.example.com
```

Before starting, point the domain's DNS record at the server and allow TCP ports 80 and 443 through the firewall. Then run:

```sh
podman compose --profile production up -d --build
podman compose --profile production ps
podman compose --profile production logs -f caddy
```

Caddy obtains and renews the TLS certificate automatically. For a server that uses rootless Podman, ports 80 and 443 may require additional setup; otherwise run the commands with `sudo`.

To deploy an update:

```sh
git pull
podman compose --profile production up -d --build
```

## Native development

Start PostgreSQL and apply the schema:

```sh
createdb urlshortener
psql urlshortener -f backend/migrations/001_create_urls.sql
```

Start the API in one terminal:

```sh
cd backend
DATABASE_URL='postgres://urlapp:password@localhost:5432/urlshortener?sslmode=disable' go run ./cmd/api
```

Install frontend dependencies and start Vite in another terminal:

```sh
cd frontend
npm install
VITE_SHORT_URL_BASE=http://localhost:4000 npm run dev
```

Open the Vite URL shown in the terminal (normally <http://localhost:5173>). Vite proxies `/api` and short-link requests to the API on port 4000.

Run backend tests with:

```sh
cd backend
go test ./...
```

## Troubleshooting

- **`Could not create the short URL`**: inspect `podman compose logs api db`. Confirm PostgreSQL is healthy and that the migration has been applied.
- **Schema changes are not applied**: migrations in `/docker-entrypoint-initdb.d` run only for a new PostgreSQL volume. Apply the migration manually with `psql`, or recreate the volume only when its data can be discarded.
- **Port already in use**: change `APP_PORT` in `.env`, then recreate the web service with `podman compose up -d --build`.
- **Production certificate failure**: verify DNS resolves to the server and that ports 80/443 are reachable from the internet.
