# 🚀 Goli CI/CD

![Goli](img/GOLI.jpg)

**Goli** is a lightweight, self-hosted CI/CD deployment tool designed for developers who want to deploy applications on their own servers without relying on cloud platforms like Azure, AWS, or GitHub Actions. It provides a simple REST API and beautiful web UI to manage deployment pipelines, execute Docker operations, and automate your deployment workflows.

## ✨ Features

- **🎯 Pipeline Management**: Define deployment pipelines using YAML
- **🐳 Docker Integration**: Full Docker support (pull, run, start, stop, manage containers)
- **📊 Web UI**: Modern, intuitive interface for managing pipelines and jobs
- **👥 User Management**: Multi-user support with admin/user roles
- **📝 Real-time Logs**: Live job execution logs with WebSocket support
- **🔄 Job Queue**: Asynchronous job processing with retry support
- **🔐 Secure**: API key authentication and secure user management
- **⚡ Lightweight**: Minimal resource footprint, runs on any Linux server

## 📋 Prerequisites

- **Operating System**: Ubuntu/Debian Linux (or any Linux distribution)
- **Docker**: Docker must be installed and running on your server
- **Root Access**: Required for installation (to create system user and service)
- **Go**: Will be automatically installed if not present (Go 1.20.3+)

## 🚀 Installation

### Step 1: Clone or Download Goli

```bash
git clone <repository-url>
cd Goli
```

### Step 2: Run the Installation Script

The installation script will:
- Create a system user `goli` for secure service execution
- Install Go (if not already installed)
- Compile and install the Goli binary
- Set up systemd service
- Generate a secure authentication key
- Start the Goli service

```bash
sudo chmod +x install.sh
sudo ./install.sh
```

Select option `1` to install Goli.

### Step 3: Complete Initial Setup

After installation, access the Goli UI:

1. Open your browser and navigate to: `http://your-server-ip:8125`
2. Complete the initial setup wizard:
   - Create your admin user (default username: `goli`)
   - Configure application settings
   - Your temporary auth key will be displayed in the terminal

### Step 4: Configure Your Server

#### Firewall Configuration

If you have a firewall enabled, allow the Goli port:

```bash
sudo ufw allow 8125/tcp
```

#### Reverse Proxy (Optional)

For production, you may want to set up a reverse proxy (Nginx/Apache) to:
- Use HTTPS/SSL
- Use a custom domain
- Hide the port number

## 🔗 Connecting Goli with GitHub

Goli can be triggered from GitHub in several ways:

### Method 1: GitHub Actions Workflow

Create a GitHub Actions workflow that triggers Goli when code is pushed:

1. **Create a workflow file** in your repository: `.github/workflows/deploy.yml`

```yaml
name: Deploy to Server

on:
  push:
    branches:
      - main
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Trigger Goli Pipeline
        run: |
          curl -X POST \
            -H "Authorization: Goli-Auth-Key YOUR_AUTH_KEY" \
            -H "Content-Type: application/json" \
            -d '{
              "name": "Deploy from GitHub - ${{ github.sha }}",
              "triggered_by": "GitHub Actions"
            }' \
            http://your-server-ip:8125/api/v1/pipelines/PIPELINE_ID/run
```

2. **Get your Auth Key**: Find it in `/goli/config/config.toml` or in the Goli UI Settings
3. **Get Pipeline ID**: Create a pipeline in Goli UI, note the ID
4. **Add as GitHub Secret**: Store your auth key as a GitHub secret for security

### Method 2: GitHub Webhook (Direct Integration)

For direct webhook integration, you can set up a webhook receiver:

1. **Create a webhook endpoint** in your application that calls Goli
2. **Configure GitHub Webhook**:
   - Go to your repository → Settings → Webhooks
   - Add webhook URL: `http://your-server-ip:8125/api/v1/pipelines/PIPELINE_ID/run`
   - Content type: `application/json`
   - Secret: Your Goli auth key
   - Events: Select "Push" or "Just the push event"

### Method 3: Manual Trigger via API

You can trigger pipelines manually using curl or any HTTP client:

```bash
curl -X POST \
  -H "Authorization: Goli-Auth-Key YOUR_AUTH_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Manual Deployment",
    "triggered_by": "Manual"
  }' \
  http://your-server-ip:8125/api/v1/pipelines/1/run
```

## 📝 Creating and Managing Pipelines

### Pipeline Definition Format

Pipelines are defined in YAML format. Here's an example:

```yaml
name: "Deploy Web Application"
description: "Deploy a web application using Docker"
steps:
  - name: "Pull Latest Image"
    type: "docker"
    action: "pull"
    config:
      image: "myapp:latest"
    retry: 2
    on_failure: "stop"

  - name: "Stop Old Container"
    type: "docker"
    action: "stop"
    config:
      container: "myapp-container"
    on_failure: "continue"

  - name: "Remove Old Container"
    type: "docker"
    action: "rm"
    config:
      container: "myapp-container"
    on_failure: "continue"

  - name: "Run New Container"
    type: "docker"
    action: "run"
    config:
      container: "myapp-container"
      image: "myapp:latest"
      ports: ["8080:80"]
      env:
        NODE_ENV: "production"
      volumes: ["/app/data:/data"]
    retry: 1
    on_failure: "stop"

  - name: "Health Check"
    type: "shell"
    action: "run"
    config:
      command: "curl"
      args: ["-f", "http://localhost:8080/health"]
    on_failure: "stop"
```

### Step Types

#### Docker Steps

**Actions:**
- `pull`: Pull a Docker image
- `run`: Run a Docker container
- `start`: Start an existing container
- `stop`: Stop a running container
- `rm`: Remove a container

**Example Docker Run:**
```yaml
- name: "Run Application"
  type: "docker"
  action: "run"
  config:
    container: "myapp"
    image: "myapp:latest"
    ports: ["8080:80", "443:443"]
    env:
      DATABASE_URL: "postgres://..."
      API_KEY: "secret-key"
    volumes: ["/host/path:/container/path"]
    cmd: ["npm", "start"]
```

#### Shell Steps

Execute shell commands:

```yaml
- name: "Run Tests"
  type: "shell"
  action: "run"
  config:
    command: "npm"
    args: ["test"]
```

#### Script Steps

Execute custom scripts:

```yaml
- name: "Custom Script"
  type: "script"
  action: "run"
  config:
    script: |
      #!/bin/bash
      echo "Deploying..."
      ./deploy.sh
    shell: "bash"
```

### Step Options

- **retry**: Number of retry attempts (default: 1)
- **on_failure**: What to do on failure
  - `stop`: Stop pipeline execution (default)
  - `continue`: Continue to next step

### Creating Pipelines via UI

1. Navigate to **Pipelines** tab in the Goli UI
2. Click **+ Upload Pipeline**
3. Select your YAML file
4. Optionally provide a name and description
5. Click **Upload**

### Running Pipelines

**Via UI:**
1. Go to Pipelines tab
2. Click **Run** next to the pipeline
3. Monitor execution in the Jobs tab

**Via API:**
```bash
POST /api/v1/pipelines/{id}/run
```

## 🏗️ Architecture & Code Structure

### Project Structure

```
Goli/
├── goli/                    # Backend (Go)
│   ├── main.go             # Application entry point
│   ├── handler/            # HTTP request handlers
│   │   ├── jobs.go        # Job management
│   │   ├── pipelines.go   # Pipeline management
│   │   ├── config.go      # Configuration management
│   │   ├── users.go       # User management
│   │   └── triggers.go    # Docker operations
│   ├── database/          # Database layer
│   │   ├── database.go    # Database initialization
│   │   ├── job_repository.go
│   │   ├── pipeline_repository.go
│   │   └── user_repository.go
│   ├── models/            # Data models
│   │   ├── job.go
│   │   ├── pipeline.go
│   │   └── user.go
│   ├── pipeline/          # Pipeline execution engine
│   │   ├── parser.go      # YAML parsing
│   │   ├── executor.go    # Step execution
│   │   └── validator.go   # Pipeline validation
│   ├── queue/             # Job queue system
│   │   └── job_queue.go   # Worker queue implementation
│   ├── websocket/         # Real-time updates
│   │   └── hub.go         # WebSocket hub
│   ├── middlewares/       # HTTP middlewares
│   │   └── auth.go        # Authentication
│   └── auxiliary/         # Utilities
│       └── toml.go        # Config file management
├── frontend/              # Frontend (Vue.js)
│   └── src/
│       ├── components/    # Vue components
│       ├── api/          # API client
│       └── App.vue       # Main app component
├── install.sh            # Installation script
└── utils/                # Utility files
    ├── config.toml       # Config template
    └── goli.service      # Systemd service file
```

### Key Components

#### 1. **Job Queue System** (`goli/queue/`)
- Asynchronous job processing
- Worker pool pattern (3 workers by default)
- Automatic retry mechanism
- Real-time status updates via WebSocket

#### 2. **Pipeline Executor** (`goli/pipeline/`)
- Parses YAML pipeline definitions
- Executes steps sequentially
- Supports Docker, Shell, and Script steps
- Handles step failures and retries

#### 3. **Database Layer** (`goli/database/`)
- SQLite database for persistence
- Stores pipelines, jobs, job steps, and users
- Repository pattern for data access

#### 4. **WebSocket Hub** (`goli/websocket/`)
- Real-time job status updates
- Live log streaming
- Broadcasts to all connected clients

#### 5. **Authentication** (`goli/middlewares/`)
- API key-based authentication
- User management with roles (admin/user)
- Secure password hashing (bcrypt)

## 🔌 API Endpoints

### Authentication

All API requests require the `Authorization` header:
```
Authorization: Goli-Auth-Key YOUR_AUTH_KEY
```

### Pipeline Endpoints

```
GET    /api/v1/pipelines              # List all pipelines
POST   /api/v1/pipelines              # Create pipeline (JSON)
POST   /api/v1/pipelines/upload       # Upload pipeline (YAML file)
GET    /api/v1/pipelines/{id}         # Get pipeline details
POST   /api/v1/pipelines/{id}/run     # Run a pipeline
```

### Job Endpoints

```
GET    /api/v1/jobs                   # List all jobs
POST   /api/v1/jobs                   # Create a job
GET    /api/v1/jobs/{id}              # Get job details with logs
```

### Configuration Endpoints

```
GET    /api/v1/config                 # Get configuration
POST   /api/v1/config                 # Update configuration
```

### User Management Endpoints

```
GET    /api/v1/users                  # List all users
POST   /api/v1/users                  # Create user
PUT    /api/v1/users/{id}            # Update user
DELETE /api/v1/users/{id}            # Delete user
```

### Docker Endpoints

```
POST   /api/v1/docker/container/start
POST   /api/v1/docker/container/stop
POST   /api/v1/docker/container/rm
POST   /api/v1/docker/container/run
POST   /api/v1/docker/image/pull
POST   /api/v1/docker/image/rm
POST   /api/v1/docker/ps
POST   /api/v1/docker/images
```

### WebSocket

```
WS     /ws                            # WebSocket connection for real-time updates
```

## 📖 Example Use Cases

### Use Case 1: Deploy Node.js Application

```yaml
name: "Deploy Node.js App"
steps:
  - name: "Pull Image"
    type: "docker"
    action: "pull"
    config:
      image: "node:18-alpine"
  
  - name: "Stop Old Container"
    type: "docker"
    action: "stop"
    config:
      container: "my-node-app"
    on_failure: "continue"
  
  - name: "Run New Container"
    type: "docker"
    action: "run"
    config:
      container: "my-node-app"
      image: "node:18-alpine"
      ports: ["3000:3000"]
      env:
        NODE_ENV: "production"
```

### Use Case 2: Database Migration

```yaml
name: "Database Migration"
steps:
  - name: "Run Migrations"
    type: "script"
    action: "run"
    config:
      script: |
        cd /app
        npm run migrate
      shell: "bash"
    retry: 2
```

### Use Case 3: Multi-Step Deployment

```yaml
name: "Full Stack Deployment"
steps:
  - name: "Pull Backend Image"
    type: "docker"
    action: "pull"
    config:
      image: "backend:latest"
  
  - name: "Pull Frontend Image"
    type: "docker"
    action: "pull"
    config:
      image: "frontend:latest"
  
  - name: "Deploy Backend"
    type: "docker"
    action: "run"
    config:
      container: "backend"
      image: "backend:latest"
      ports: ["8000:8000"]
  
  - name: "Deploy Frontend"
    type: "docker"
    action: "run"
    config:
      container: "frontend"
      image: "frontend:latest"
      ports: ["80:80"]
```

## 🔧 Configuration

Configuration is stored in `/goli/config/config.toml`:

```toml
[constants]
auth_key = "your-secure-auth-key"
port = "8125"
setup_complete = true
```

You can update configuration via:
- **UI**: Settings → General Settings
- **API**: `POST /api/v1/config`
- **File**: Directly edit `/goli/config/config.toml`

## 🛠️ Service Management

### Start/Stop Service

```bash
sudo systemctl start goli.service
sudo systemctl stop goli.service
sudo systemctl restart goli.service
sudo systemctl status goli.service
```

### View Logs

```bash
sudo journalctl -u goli.service -f
```

### Update Goli

```bash
sudo ./install.sh
# Select option 1 (Update)
```

### Uninstall

```bash
sudo ./install.sh
# Select option 2 (Remove)
```

## 🔒 Security Best Practices

1. **Change Default Auth Key**: Update the auth key after installation
2. **Use HTTPS**: Set up reverse proxy with SSL certificate
3. **Firewall**: Only expose necessary ports
4. **User Management**: Create separate users for different team members
5. **Regular Updates**: Keep Goli updated to latest version
6. **Backup**: Regularly backup `/goli/data/goli.db` database

## 🐛 Troubleshooting

### Service Won't Start

```bash
# Check service status
sudo systemctl status goli.service

# Check logs
sudo journalctl -u goli.service -n 50

# Verify permissions
ls -la /goli/
ls -la /usr/local/sbin/goli/
```

### Can't Access UI

1. Check if service is running: `sudo systemctl status goli.service`
2. Verify port is open: `sudo netstat -tlnp | grep 8125`
3. Check firewall: `sudo ufw status`
4. Verify config: Check `/goli/config/config.toml`

### Pipeline Execution Fails

1. Check job logs in the UI
2. Verify Docker is running: `sudo systemctl status docker`
3. Check Docker permissions for `goli` user
4. Review pipeline YAML syntax

## 📚 Additional Resources

- **UI Documentation**: See `README-UI.md` for frontend details
- **Example Pipelines**: Check `goli/web/example-pipeline.yaml`
- **API Documentation**: All endpoints are documented in code comments

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

[Your License Here]

## 🙏 Acknowledgments

Goli was created to provide a simple, self-hosted alternative to cloud-based CI/CD solutions, giving developers full control over their deployment infrastructure.

---

**Need Help?** Open an issue on GitHub or check the troubleshooting section above.

**Happy Deploying! 🚀**
