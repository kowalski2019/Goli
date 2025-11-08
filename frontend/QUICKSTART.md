# Goli Frontend - Quick Start Guide

## 🚀 Getting Started

### 1. Install Dependencies

```bash
cd goli/frontend
npm install
```

### 2. Start Development Server

```bash
npm run dev
```

The frontend will run on `http://localhost:5173` (or another port if 5173 is taken).

The Vite dev server is configured to proxy API requests to the Go backend at `http://localhost:8125`.

### 3. Start the Go Backend

In a separate terminal:

```bash
cd goli
go run main.go
```

The backend should be running on `http://localhost:8125`.

## 📦 Building for Production

```bash
npm run build
```

This will build the Vue app and output the files to `goli/web/` directory, which is served by the Go backend.

## ✨ Features

### 📊 Dashboard
- Real-time job statistics
- Recent jobs overview
- Quick access to job details

### 📋 Jobs Management
- View all jobs with status
- Create new jobs
- View job details with step-by-step progress
- **View logs for each step** - Click on any job and select "Logs" to see detailed step-by-step logs
- Real-time job status updates

### 🔄 Pipeline Management
- Upload YAML pipeline definitions
- Run pipelines with one click
- View pipeline details

### 📜 Logs Viewer (NEW!)
- **Step-by-step log viewing** - Navigate through pipeline steps
- **Real-time log updates** - Logs refresh automatically for running jobs
- **Side-by-side navigation** - Steps list on the left, logs on the right
- **Syntax highlighting** - Terminal-style log display
- **Error highlighting** - Errors are clearly marked
- **Job-level logs** - View overall job logs at the bottom

## 🎨 UI Features

- Modern, responsive design with Tailwind CSS
- Real-time WebSocket updates
- Modal dialogs for detailed views
- Status badges with color coding
- Loading states and error handling
- Smooth animations and transitions

## 📁 Project Structure

```
frontend/
├── src/
│   ├── api/
│   │   └── client.js          # API client & WebSocket
│   ├── components/
│   │   ├── Dashboard.vue      # Main dashboard
│   │   ├── Jobs.vue           # Jobs list
│   │   ├── Pipelines.vue      # Pipelines list
│   │   ├── JobDetailsModal.vue # Job details
│   │   ├── JobLogsModal.vue   # Step-by-step logs viewer ⭐
│   │   ├── CreateJobModal.vue  # Create job form
│   │   └── UploadPipelineModal.vue # Upload pipeline
│   ├── App.vue                # Main app
│   ├── main.js                # Entry point
│   └── style.css              # Tailwind imports
├── package.json
├── vite.config.js
└── tailwind.config.js
```

## 🔍 How to View Logs

1. **From Jobs List:**
   - Click on any job row
   - Click "Logs" button
   - View step-by-step logs in the modal

2. **From Job Details:**
   - Click on a job to see details
   - Click "View All Logs" button
   - Navigate through steps on the left sidebar
   - View logs for each step on the right

3. **Real-time Updates:**
   - Logs automatically refresh every 2 seconds for running jobs
   - WebSocket updates notify you of status changes

## 🛠️ Development Tips

- The dev server has hot module replacement (HMR) - changes appear instantly
- API calls are proxied to the backend automatically
- WebSocket connections are handled automatically
- Check browser console for debugging info

## 📝 Next Steps

- Add more visualizations (charts, graphs)
- Add pipeline editor
- Add user authentication UI
- Add settings page
- Add job history and filtering

