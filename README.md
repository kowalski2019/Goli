# 🚀 Goli CI/CD

![Goli](img/GOLI.jpg)

**Goli** is a lightweight, self-hosted CI/CD platform that gives you full control over your deployment infrastructure. Deploy applications on your own servers without relying on cloud platforms.

## ✨ Features

- **🎯 Pipeline Management**: Define deployment pipelines using YAML
- **🐳 Docker Integration**: Full Docker support (containers, images, compose)
- **📊 Modern Web UI**: Beautiful Vue.js interface with real-time updates
- **👥 User Management**: Multi-user support with roles and 2FA (Email/SMS)
- **📝 Real-time Logs**: Live job execution logs with WebSocket support
- **🔐 Secure**: Session-based authentication with Bearer tokens
- **⚡ Lightweight**: Minimal resource footprint, runs on any Linux server

## 🚀 Quick Start

### Installation

```bash
git clone <repository-url>
cd Goli
sudo chmod +x install.sh
sudo ./install.sh
```

Select option `1` to install. The script will:
- Install Go (if needed)
- Compile and install Goli
- Set up systemd service
- Generate authentication keys

### First Steps

1. **Access the UI**: Open `http://your-server-ip:8125`
2. **Complete Setup**: Use the setup wizard to create your admin account
3. **Create Pipeline**: Upload a YAML pipeline or create one via UI
4. **Run Jobs**: Execute pipelines and monitor in real-time

## 📚 Documentation

- **[Installation Guide](INSTALLATION.md)** - Detailed installation and configuration
- **[API Reference](API.md)** - Complete API documentation
- **[Pipeline Guide](PIPELINES.md)** - Creating and managing pipelines
- **[UI Documentation](README-UI.md)** - Frontend development guide

## 🏗️ Architecture

```
Goli/
├── goli/              # Backend (Go)
│   ├── handler/       # HTTP handlers
│   ├── pipeline/      # Pipeline execution engine
│   ├── queue/         # Job queue system
│   └── database/      # SQLite persistence
├── frontend/          # Frontend (Vue.js + Tailwind CSS)
└── utils/             # Configuration templates
```

**Key Components:**
- **Job Queue**: Asynchronous processing with worker pool
- **Pipeline Executor**: Supports Docker, Shell, and Script steps
- **WebSocket Hub**: Real-time updates and log streaming
- **Authentication**: Session-based with 2FA support

## 🔌 API Overview

Goli provides a RESTful API with WebSocket support:

- **Authentication**: `POST /api/v1/auth/login` (Bearer token)
- **Pipelines**: Create, list, run pipelines
- **Jobs**: Monitor and manage job execution
- **Docker**: Direct Docker operations
- **WebSocket**: Real-time updates at `/ws`

See [API.md](API.md) for complete documentation.

## 📝 Pipeline Example

```yaml
name: "Deploy Application"
steps:
  - name: "Pull Image"
    type: "docker"
    action: "pull"
    config:
      image: "myapp:latest"
  
  - name: "Run Container"
    type: "docker"
    action: "run"
    config:
      container: "myapp"
      image: "myapp:latest"
      ports: ["8080:80"]
```

See [PIPELINES.md](PIPELINES.md) for detailed pipeline documentation.

## 🛠️ Service Management

```bash
# Start/Stop/Restart
sudo systemctl start goli.service
sudo systemctl stop goli.service
sudo systemctl restart goli.service

# View logs
sudo journalctl -u goli.service -f

# Update
sudo ./install.sh  # Select option 1
```

## 🔒 Security

- Session-based authentication with Bearer tokens
- 2FA support (Email/SMS via Twilio)
- Secure password hashing (bcrypt)
- API key support for automation
- User roles (admin/user)

## 📖 Example Use Cases

- **Deploy Docker Applications**: Pull images, manage containers
- **Database Migrations**: Run scripts and migrations
- **Multi-Service Deployments**: Orchestrate complex deployments
- **CI/CD Integration**: Trigger from GitHub Actions, webhooks

## 🐛 Troubleshooting

**Service won't start?**
```bash
sudo systemctl status goli.service
sudo journalctl -u goli.service -n 50
```

**Can't access UI?**
- Check service status
- Verify port 8125 is open
- Check firewall settings

**Pipeline fails?**
- Check job logs in UI
- Verify Docker is running
- Review pipeline YAML syntax

## 📄 License

MIT License

## 🤝 Contributing

Contributions welcome! Please feel free to submit a Pull Request.

---

**Need Help?** Check the detailed documentation files or open an issue.

**Happy Deploying! 🚀**
