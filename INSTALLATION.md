# Installation Guide

Complete installation and configuration guide for Goli CI/CD.

## Prerequisites

- **Operating System**: Ubuntu/Debian Linux (or any Linux distribution)
- **Docker**: Docker must be installed and running
- **Root Access**: Required for installation
- **Go**: Will be automatically installed if not present (Go 1.20.3+)

## Installation

### Step 1: Clone Repository

```bash
git clone <repository-url>
cd Goli
```

### Step 2: Run Installation Script

```bash
sudo chmod +x install.sh
sudo ./install.sh
```

Select option `1` to install Goli.

The script will:
- Create system user `goli` for secure service execution
- Install Go (if not already installed)
- Compile and install the Goli binary
- Set up systemd service
- Generate secure authentication key
- Start the Goli service

### Step 3: Complete Initial Setup

1. **Access the UI**: Open `http://your-server-ip:8125`
2. **Setup Wizard**: Complete the initial setup:
   - Enter setup password (shown in terminal during installation)
   - Create admin user (default username: `goli`)
   - Configure system settings (port, auth key)
3. **Login**: Use your admin credentials to log in

## Configuration

### Firewall Setup

If you have a firewall enabled:

```bash
sudo ufw allow 8125/tcp
```

### Reverse Proxy (Optional)

For production, set up a reverse proxy for HTTPS:

**Nginx Example:**
```nginx
server {
    listen 80;
    server_name goli.example.com;

    location / {
        proxy_pass http://localhost:8125;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### Configuration File

Main configuration is stored in `/goli/config/config.toml`:

```toml
[constants]
auth_key = "your-auth-key"
port = "8125"
setup_complete = true

# GitHub Integration (Optional)
gh_username = "your-username"
gh_access_token = "your-token"

# SMTP Configuration (Optional)
smtp_host = "smtp.example.com"
smtp_port = "587"
smtp_user = "user@example.com"
smtp_pass = "password"
smtp_from = "noreply@example.com"
smtp_from_name = "Goli CI/CD"

# Twilio for SMS 2FA (Optional)
twilio_sid = "your-sid"
twilio_token = "your-token"
twilio_from = "+1234567890"
```

**Update via UI:**
- Go to Settings → System Configuration
- Update values and click "Update Configuration"

**Update via API:**
```bash
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"port": "8125", "smtp_host": "smtp.example.com"}' \
  http://localhost:8125/api/v1/config
```

## Service Management

### Start/Stop/Restart

```bash
sudo systemctl start goli.service
sudo systemctl stop goli.service
sudo systemctl restart goli.service
sudo systemctl status goli.service
```

### View Logs

```bash
# Follow logs
sudo journalctl -u goli.service -f

# Last 50 lines
sudo journalctl -u goli.service -n 50

# Since today
sudo journalctl -u goli.service --since today
```

### Enable Auto-Start

```bash
sudo systemctl enable goli.service
```

## Updating Goli

```bash
cd Goli
git pull
sudo ./install.sh
# Select option 1 (Update)
```

## Uninstallation

```bash
sudo ./install.sh
# Select option 2 (Remove)
```

This will:
- Stop the service
- Remove the binary
- Remove systemd service
- Keep data and config (for safety)

To completely remove:
```bash
sudo rm -rf /goli
sudo rm -rf /usr/local/sbin/goli
```

## Docker Setup

Goli requires Docker to be installed and the `goli` user to have Docker access:

```bash
# Add goli user to docker group
sudo usermod -aG docker goli

# Verify Docker access
sudo -u goli docker ps
```

## Troubleshooting

### Service Won't Start

```bash
# Check status
sudo systemctl status goli.service

# Check logs
sudo journalctl -u goli.service -n 50

# Verify permissions
ls -la /goli/
ls -la /usr/local/sbin/goli
```

### Can't Access UI

1. **Check service**: `sudo systemctl status goli.service`
2. **Check port**: `sudo netstat -tlnp | grep 8125`
3. **Check firewall**: `sudo ufw status`
4. **Check config**: Verify `/goli/config/config.toml`

### Permission Issues

```bash
# Fix ownership
sudo chown -R goli:goli /goli

# Fix Docker permissions
sudo usermod -aG docker goli
```

### Database Issues

Database is stored at `/goli/data/goli.db`. To reset:

```bash
sudo systemctl stop goli.service
sudo mv /goli/data/goli.db /goli/data/goli.db.backup
sudo systemctl start goli.service
```

## Backup

**Important files to backup:**
- `/goli/data/goli.db` - Database
- `/goli/config/config.toml` - Configuration

```bash
# Backup script
sudo tar -czf goli-backup-$(date +%Y%m%d).tar.gz \
  /goli/data/goli.db \
  /goli/config/config.toml
```

## Security Considerations

1. **Change default auth key** after installation
2. **Use HTTPS** via reverse proxy in production
3. **Firewall rules** - only expose necessary ports
4. **Regular updates** - keep Goli updated
5. **User management** - create separate users for team members
6. **2FA** - enable 2FA for admin accounts

## Next Steps

After installation:
1. Complete setup wizard
2. Create your first pipeline
3. Configure SMTP (for email notifications)
4. Set up reverse proxy (for HTTPS)
5. Review security settings

See [PIPELINES.md](PIPELINES.md) for creating pipelines and [API.md](API.md) for API usage.

