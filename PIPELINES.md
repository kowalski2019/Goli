# Pipeline Guide

Complete guide to creating and managing pipelines in Goli.

## Pipeline Structure

Pipelines are defined in YAML format with the following structure:

```yaml
name: "Pipeline Name"
description: "Optional description"
steps:
  - name: "Step Name"
    type: "docker" | "shell" | "script"
    action: "action-name"
    config:
      # Step-specific configuration
      # You can use variables: ${VAR_NAME} or {{VAR_NAME}}
    retry: 1                    # Optional: retry attempts (default: 1)
    on_failure: "stop"         # Optional: "stop" or "continue" (default: "stop")
```

## Variables and Secrets

Pipelines support variables and secrets that can be used throughout the pipeline definition. This allows you to:
- Keep sensitive values (API keys, passwords) out of your git repository
- Reuse common values across multiple pipelines
- Easily update values without editing YAML files

### Using Variables

Variables can be referenced in your pipeline YAML using either syntax:
- `${VAR_NAME}` - Standard variable syntax
- `{{VAR_NAME}}` - Alternative syntax

**Example:**
```yaml
name: "Deploy Application"
steps:
  - name: "Pull Image"
    type: "docker"
    action: "pull"
    config:
      image: "${IMAGE_NAME}:${IMAGE_TAG}"
  
  - name: "Run Container"
    type: "docker"
    action: "run"
    config:
      container: "${CONTAINER_NAME}"
      image: "${IMAGE_NAME}:${IMAGE_TAG}"
      env:
        DATABASE_URL: "${DATABASE_URL}"
        API_KEY: "${API_SECRET_KEY}"
```

Variables are substituted at execution time, so secrets are never exposed in logs or stored in git.

### Managing Variables

Variables can be managed through the UI:
1. When creating or editing a pipeline, use the **Variables & Secrets** section
2. Click **Add Variable** to create a new variable
3. Enter the variable name (e.g., `DATABASE_URL`)
4. Enter the value
5. Check **Mark as secret** to mask the value (it will be hidden in the UI)
6. Variables marked as secrets are preserved when editing (you must enter a new value to update them)

**Best Practices:**
- Use secrets for sensitive data (passwords, API keys, tokens)
- Use regular variables for non-sensitive configuration (image names, ports, etc.)
- Use descriptive variable names (e.g., `DATABASE_URL` instead of `DB`)
- Keep variable names in UPPERCASE for consistency

## Step Types

### Docker Steps

Execute Docker operations.

**Actions:**
- `pull`: Pull a Docker image
- `run`: Run a new container
- `start`: Start an existing container
- `stop`: Stop a running container
- `rm`: Remove a container

**Pull Image:**
```yaml
- name: "Pull Latest Image"
  type: "docker"
  action: "pull"
  config:
    image: "myapp:latest"
```

**Run Container:**
```yaml
- name: "Run Application"
  type: "docker"
  action: "run"
  config:
    container: "myapp-container"
    image: "myapp:latest"
    ports: ["8080:80", "443:443"]
    env:
      NODE_ENV: "production"
      DATABASE_URL: "postgres://..."
    volumes: ["/host/path:/container/path"]
    cmd: ["npm", "start"]          # Optional: command to run
    network: "my-network"          # Optional: network name
```

**Container Operations:**
```yaml
- name: "Stop Container"
  type: "docker"
  action: "stop"
  config:
    container: "myapp-container"
  on_failure: "continue"           # Continue even if container doesn't exist
```

### Shell Steps

Execute shell commands.

```yaml
- name: "Run Tests"
  type: "shell"
  action: "run"
  config:
    command: "npm"
    args: ["test"]
```

**Example:**
```yaml
- name: "Build Application"
  type: "shell"
  action: "run"
  config:
    command: "make"
    args: ["build"]
  retry: 2
```

### Script Steps

Execute custom scripts.

```yaml
- name: "Custom Deployment"
  type: "script"
  action: "run"
  config:
    script: |
      #!/bin/bash
      echo "Deploying..."
      ./deploy.sh
      echo "Deployment complete"
    shell: "bash"                  # Optional: shell to use (default: "sh")
```

## Step Options

### Retry

Number of retry attempts if step fails:

```yaml
- name: "Unreliable Step"
  type: "shell"
  action: "run"
  config:
    command: "curl"
    args: ["-f", "http://example.com"]
  retry: 3                         # Will retry up to 3 times
```

### On Failure

Control pipeline behavior when a step fails:

```yaml
- name: "Optional Step"
  type: "docker"
  action: "stop"
  config:
    container: "old-container"
  on_failure: "continue"           # Continue to next step even if this fails
```

Options:
- `stop`: Stop pipeline execution (default)
- `continue`: Continue to next step

## Complete Examples

### Example 1: Deploy Node.js Application

```yaml
name: "Deploy Node.js App"
description: "Deploy a Node.js application using Docker"
steps:
  - name: "Pull Latest Image"
    type: "docker"
    action: "pull"
    config:
      image: "myapp:latest"
    retry: 2

  - name: "Stop Old Container"
    type: "docker"
    action: "stop"
    config:
      container: "myapp"
    on_failure: "continue"

  - name: "Remove Old Container"
    type: "docker"
    action: "rm"
    config:
      container: "myapp"
    on_failure: "continue"

  - name: "Run New Container"
    type: "docker"
    action: "run"
    config:
      container: "myapp"
      image: "myapp:latest"
      ports: ["3000:3000"]
      env:
        NODE_ENV: "production"
        PORT: "3000"
      volumes: ["/app/data:/data"]
    retry: 1

  - name: "Health Check"
    type: "shell"
    action: "run"
    config:
      command: "curl"
      args: ["-f", "http://localhost:3000/health"]
    retry: 3
```

### Example 2: Database Migration

```yaml
name: "Database Migration"
steps:
  - name: "Run Migrations"
    type: "script"
    action: "run"
    config:
      script: |
        cd /app
        export DATABASE_URL="postgres://user:pass@localhost/db"
        npm run migrate
      shell: "bash"
    retry: 2
    on_failure: "stop"
```

### Example 3: Multi-Service Deployment

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
      env:
        DATABASE_URL: "postgres://..."
      network: "app-network"

  - name: "Deploy Frontend"
    type: "docker"
    action: "run"
    config:
      container: "frontend"
      image: "frontend:latest"
      ports: ["80:80"]
      network: "app-network"
```

## Creating Pipelines

### Via UI Editor (Recommended)

The Goli UI provides a full-featured code editor with syntax highlighting and indentation support:

1. Navigate to **Pipelines** tab
2. Click **Create Pipeline** button
3. Enter pipeline name and description
4. Write or paste your YAML definition in the code editor
   - The editor provides YAML syntax highlighting
   - Auto-indentation and code formatting
   - Dark/light theme toggle
5. Add variables and secrets in the **Variables & Secrets** section (optional)
6. Click **Create Pipeline** to save

**Editing Pipelines:**
1. Navigate to **Pipelines** tab
2. Click **Edit** next to the pipeline you want to modify
3. Make your changes in the full-page editor
4. Update variables if needed (secrets are preserved unless you enter a new value)
5. Click **Save Changes**

### Via UI Upload

1. Navigate to **Pipelines** tab
2. Click **Upload YAML** button
3. Select your YAML file
4. Optionally provide name and description
5. Click **Upload & Create**

**Note:** After uploading, you can edit the pipeline using the full-page editor.

### Via API

**Upload YAML File:**
```bash
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -F "yaml_file=@pipeline.yaml" \
  -F "name=My Pipeline" \
  http://your-server:8125/api/v1/pipelines/upload
```

**Create from JSON:**
```bash
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Pipeline",
    "definition": "yaml_content_here"
  }' \
  http://your-server:8125/api/v1/pipelines
```

## Running Pipelines

### Via UI

1. Go to **Pipelines** tab
2. Click **Run** next to the pipeline
3. Monitor execution in **Jobs** tab

### Via API

```bash
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deployment Run",
    "triggered_by": "Manual"
  }' \
  http://your-server:8125/api/v1/pipelines/1/run
```

## Best Practices

1. **Use descriptive names**: Clear step names help with debugging
2. **Handle failures gracefully**: Use `on_failure: "continue"` for optional steps
3. **Add retries**: For network operations, add retry logic
4. **Health checks**: Always verify deployment success
5. **Use Variables & Secrets**: Store sensitive data as secrets, not in YAML
6. **Idempotency**: Design steps to be safely re-runnable
7. **Version control**: Keep pipeline YAML in git, but exclude secrets
8. **Documentation**: Add descriptions to pipelines and steps for clarity

## Troubleshooting

**Pipeline fails immediately?**
- Check YAML syntax
- Verify all required config fields are present
- Check Docker is running

**Step keeps failing?**
- Review step logs in the UI
- Verify Docker permissions
- Check network connectivity
- Review environment variables

**Container already exists?**
- Add stop/rm steps before run
- Use `on_failure: "continue"` for cleanup steps

