# Go VPN Bot

Telegram bot for administering a VPN service.

## Requirements

- Docker
- Docker Compose
- Git
- `curl`
- `Xray-core`
- systemd

---

# Full Installation

**1. Install Docker**
```bash
curl -fsSL https://get.docker.com | sh
```

**2. Run Docker without sudo**
```bash
sudo usermod -aG docker $USER
newgrp docker
```

**3. Install Xray-core**
```bash
sudo bash -c "$(curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh)" @ install
```

**4. Clone the project**
```bash
cd ~
git clone https://github.com/tsvetkovdm/go-vpnbot.git
cd go-vpnbot
```

**5. Configure the environment**
```bash
cp .env.example .env
nano .env
```

**6. Start PostgreSQL**
```bash
mkdir -p out
make env-up
```

**7. Run database migrations**
```bash
make migrate-up
```

**8. Deploy the bot**
```bash
make vpnbot-deploy
```

**9. Verify the installation**
```bash
make ps
docker compose logs --tail=100 vpnbot
sudo systemctl status xray
```

---

## Useful Commands


**Check docker service status**
```bash
make ps
```

**View bot logs**
```bash
docker compose logs -f vpnbot
```

**View PostgreSQL logs**
```bash
docker compose logs -f vpnbot-postgres
```

**Check Xray status**
```bash
sudo systemctl status xray
sudo journalctl -u xray -f
```

**Rebuild and deploy the bot**
```bash
make vpnbot-deploy
```

**Stop the bot**
```bash
make vpnbot-undeploy
```

**Run database migrations**
```bash
make migrate-up
```

**Roll back the last migration**
```bash
make migrate-down
```

---

# Data & Files

PostgreSQL data is stored in:

```text
out/.pgdata/
```

Bot logs:

```text
out/logs/
```

Database migrations:

```text
migrations/
```

Environment configuration:

```text
.env
```

---

## Tech Stack

- Go
- Telegram Bot API
- PostgreSQL
- Docker / Docker Compose
- Xray-core
- `golang-migrate`