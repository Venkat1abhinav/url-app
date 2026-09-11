# URL App

A Go/PostgreSQL URL shortener with a React frontend. The production stack uses Caddy for automatic HTTPS and serves the UI and API from one origin, so no browser CORS setup is required.

## Production architecture

```text
Internet (80/443) -> Caddy -> web (Nginx) -> api (Go) -> db (PostgreSQL)
```

Only Caddy is publicly exposed. The web service is bound to `127.0.0.1:8080` for local diagnostics; the API and database have no host ports.

## Deploy with Podman Compose

These steps assume a Linux VPS, a domain such as `go.example.com`, and a user allowed to run Podman. Run the production stack with `sudo podman compose` unless you have explicitly configured rootless Podman to bind ports 80 and 443.

1. Install Podman and `podman-compose` using your distribution's package manager. Open TCP ports `80` and `443` in both your cloud firewall and server firewall.

2. Create a DNS `A` (and optional `AAAA`) record for your domain pointing to the server. Wait for it to resolve before starting Caddy; it needs the domain to obtain a TLS certificate.

3. Clone or copy this repository to the server:

   ```sh
   sudo mkdir -p /opt/url-app
   sudo chown "$USER" /opt/url-app
   git clone <repository-url> /opt/url-app
   cd /opt/url-app
   ```

4. Create the production environment file and secure it:

   ```sh
   cp .env.example .env
   chmod 600 .env
   openssl rand -hex 32
   ```

   Put the generated value in `POSTGRES_PASSWORD`, then set both `DOMAIN` and `VITE_SHORT_URL_BASE` to your public HTTPS URL. Example:

   ```env
   POSTGRES_PASSWORD=generated-hex-value
   POSTGRES_DB=urlshortener
   POSTGRES_USER=urlapp
   DOMAIN=go.example.com
   VITE_SHORT_URL_BASE=https://go.example.com
   APP_PORT=8080
   CORS_ALLOWED_ORIGINS=
   ```

   Use a hexadecimal password as shown above: it is safe inside the PostgreSQL connection URL assembled by Compose.

5. Build and start the full HTTPS stack:

   ```sh
   sudo podman compose --profile production up -d --build
   sudo podman compose ps
   sudo podman compose logs -f caddy
   ```

   Open `https://go.example.com`. On first startup Caddy obtains the TLS certificate and persists it in the `caddy_data` volume. PostgreSQL data persists in `postgres_data`; the schema migration runs only when that volume is first initialized.

6. Enable startup after a server reboot:

   ```sh
   sudo cp deploy/url-app.service /etc/systemd/system/url-app.service
   sudo systemctl daemon-reload
   sudo systemctl enable --now url-app
   sudo systemctl status url-app
   ```

   The systemd unit starts previously built images. Deploy new application code with the update command below.

## Updating and operating

```sh
cd /opt/url-app
git pull
sudo podman compose --profile production up -d --build
sudo podman compose ps
sudo podman compose logs -f --tail=100
```

The loopback-only health endpoint is available on the server at:

```sh
curl http://127.0.0.1:8080/health
```

Back up PostgreSQL before upgrades and keep backups off the server:

```sh
mkdir -p backups
sudo podman compose exec -T db pg_dump -U urlapp urlshortener > "backups/url-app-$(date +%F).sql"
```

To stop the stack without deleting data:

```sh
sudo podman compose --profile production down
```

Do not run `podman compose down -v` unless you intentionally want to delete the database and TLS certificate volumes.

## Local containers

For a local HTTP-only stack, copy `.env.example` to `.env`, use a local value such as `VITE_SHORT_URL_BASE=http://localhost:8080`, and do not enable the production profile:

```sh
podman compose up -d --build
curl http://127.0.0.1:8080/health
```

## Local development

Start PostgreSQL with the migration in `backend/migrations/001_create_urls.sql`, then run the API with `DATABASE_URL=postgres://... go run ./cmd/api` from `backend`. Run `VITE_SHORT_URL_BASE=http://localhost:4000 npm run dev` from `frontend`; Vite proxies API calls and short URLs to port 4000.
