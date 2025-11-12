# Goli API Reference

Complete API documentation for Goli CI/CD platform.

## Authentication

Goli supports two authentication methods:

### 1. Bearer Token (Recommended)

After login, use the session token:

```
Authorization: Bearer <session_token>
```

**Login Flow:**
1. `POST /api/v1/auth/login` - Get session token
2. If 2FA enabled, verify with `POST /api/v1/auth/2fa/verify`
3. Use returned token in `Authorization: Bearer <token>` header

### 2. API Key (Legacy/Automation)

For automation scripts and CI/CD:

```
Authorization: Goli-Auth-Key <your_auth_key>
```

Get your auth key from Settings in the UI or `/goli/config/config.toml`.

## Public Endpoints

### Setup

```
POST   /api/v1/setup/verify        # Verify setup password
GET    /api/v1/setup/status        # Check setup status
```

### Authentication

```
POST   /api/v1/auth/login          # Login (returns token or 2FA challenge)
POST   /api/v1/auth/2fa/verify     # Verify 2FA code
POST   /api/v1/auth/logout         # Logout (invalidate session)
```

## Protected Endpoints

All endpoints below require authentication.

### Pipelines

```
GET    /api/v1/pipelines              # List all pipelines
POST   /api/v1/pipelines              # Create pipeline (JSON)
POST   /api/v1/pipelines/upload       # Upload pipeline (YAML file)
GET    /api/v1/pipelines/{id}         # Get pipeline details
POST   /api/v1/pipelines/{id}/run     # Run a pipeline
```

**Create Pipeline (JSON):**
```json
{
  "name": "My Pipeline",
  "description": "Pipeline description",
  "definition": "yaml_content_here"
}
```

**Upload Pipeline (Form Data):**
- `yaml_file`: YAML file
- `name`: Optional pipeline name
- `description`: Optional description
- `run`: "true" to run immediately

**Run Pipeline:**
```json
{
  "name": "Job Name",
  "triggered_by": "Manual"
}
```

### Jobs

```
GET    /api/v1/jobs                   # List jobs (query: ?limit=50)
POST   /api/v1/jobs                   # Create a job
GET    /api/v1/jobs/{id}              # Get job details with logs
POST   /api/v1/jobs/{id}/cancel       # Cancel a running job
```

**Create Job:**
```json
{
  "name": "Job Name",
  "triggered_by": "Manual"
}
```

### Configuration

```
GET    /api/v1/config                 # Get configuration
POST   /api/v1/config                 # Update configuration
```

**Update Config:**
```json
{
  "port": "8125",
  "gh_username": "username",
  "gh_access_token": "token",
  "smtp_host": "smtp.example.com",
  "smtp_port": "587",
  "smtp_user": "user@example.com",
  "smtp_pass": "password",
  "smtp_from": "noreply@example.com",
  "smtp_from_name": "Goli CI/CD"
}
```

### User Management

```
GET    /api/v1/users                  # List all users
POST   /api/v1/users                  # Create user
PUT    /api/v1/users/{id}             # Update user
DELETE /api/v1/users/{id}             # Delete user
```

**Create User:**
```json
{
  "username": "username",
  "password": "password",
  "email": "user@example.com",
  "phone": "+1234567890",
  "role": "user"
}
```

**Update User:**
```json
{
  "email": "newemail@example.com",
  "phone": "+1234567890",
  "two_fa_email_enabled": 1,
  "two_fa_sms_enabled": 0
}
```

### Docker Operations

```
POST   /api/v1/docker/container/start    # Start container
POST   /api/v1/docker/container/stop     # Stop container
POST   /api/v1/docker/container/rm       # Remove container
POST   /api/v1/docker/container/run      # Run new container
POST   /api/v1/docker/image/pull         # Pull image
POST   /api/v1/docker/image/rm           # Remove image
POST   /api/v1/docker/ps                 # List containers
POST   /api/v1/docker/images             # List images
POST   /api/v1/docker/compose/up         # Docker Compose up
POST   /api/v1/docker/compose/down       # Docker Compose down
```

**Container Operations:**
```json
{
  "name": "container-name"
}
```

**Run Container:**
```json
{
  "name": "container-name",
  "image": "image:tag",
  "network": "network-name",
  "port_ex": "8080",
  "port_in": "80",
  "volume_ex": "/host/path",
  "volume_in": "/container/path",
  "opts": "additional-options"
}
```

**Image Operations:**
```json
{
  "image": "image:tag"
}
```

## WebSocket

```
WS     /ws                            # WebSocket connection
```

**Message Types:**
- `job_update`: Job status changed
- `log_update`: Log content updated
- `stats_update`: Statistics updated

**Connection:**
```javascript
const ws = new WebSocket('ws://your-server:8125/ws');
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  // Handle message
};
```

## Response Format

**Success:**
```json
{
  "status": "ok",
  "data": { ... }
}
```

**Error:**
```json
{
  "status": "error",
  "description": "Error message"
}
```

## Rate Limiting

Currently no rate limiting is enforced. Use responsibly.

## Error Codes

- `400`: Bad Request - Invalid input
- `401`: Unauthorized - Missing or invalid authentication
- `404`: Not Found - Resource doesn't exist
- `500`: Internal Server Error - Server error

