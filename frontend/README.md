# Goli Frontend

Vue 3 + Vite + Tailwind CSS frontend for Goli CI/CD platform.

## Features

- 🎨 Modern UI with Tailwind CSS
- 📊 Real-time dashboard with job statistics
- 📝 Pipeline management and YAML upload
- 📋 Job execution and monitoring
- 📜 Step-by-step logs viewer
- 🔄 Real-time updates via WebSocket
- 📱 Responsive design

## Development

```bash
# Install dependencies
npm install

# Start dev server
npm run dev
```

The dev server will proxy API requests to `http://localhost:8125` (the Go backend).

## Building for Production

```bash
# Build for production
npm run build
```

The build output will be in `goli/web/` directory, which is served by the Go backend.

## Project Structure

```
src/
├── api/
│   └── client.js          # API client and WebSocket setup
├── components/
│   ├── Dashboard.vue      # Main dashboard
│   ├── Jobs.vue           # Jobs list and management
│   ├── Pipelines.vue      # Pipelines list
│   ├── JobDetailsModal.vue # Job details view
│   ├── JobLogsModal.vue   # Step-by-step logs viewer
│   ├── CreateJobModal.vue # Create job form
│   └── UploadPipelineModal.vue # Upload pipeline form
├── App.vue                # Main app component
├── main.js                # App entry point
└── style.css              # Tailwind CSS imports
```

## Features in Detail

### Logs Viewer
- View logs for each step in a pipeline
- Real-time log updates for running jobs
- Side-by-side step navigation
- Syntax highlighting for log output
- Error message display

### Dashboard
- Real-time statistics (total, running, completed, failed jobs)
- Recent jobs list
- Quick access to job details and logs

### Pipeline Management
- Upload YAML pipeline definitions
- Run pipelines with one click
- View pipeline details
- Automatic validation

### Job Management
- Create jobs manually
- View job details with step-by-step progress
- Monitor job execution in real-time
- Access comprehensive logs

