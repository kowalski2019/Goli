# Goli Frontend

Modern Vue.js frontend for Goli CI/CD platform.

## Technology Stack

- **Vue 3** - Progressive JavaScript framework
- **Vite** - Next-generation frontend tooling
- **Tailwind CSS** - Utility-first CSS framework
- **WebSocket** - Real-time updates

## Features

- 🎨 Modern, responsive UI with Tailwind CSS
- 📊 Real-time dashboard with job statistics
- 📝 Pipeline management and YAML upload
- 📋 Job execution and monitoring
- 📜 Step-by-step logs viewer
- 🔄 Real-time updates via WebSocket
- 👥 User management and settings
- 🔐 Authentication with 2FA support

## Development

### Prerequisites

- Node.js 18+ and npm

### Setup

```bash
cd frontend
npm install
```

### Development Server

```bash
npm run dev
```

The dev server runs on `http://localhost:5173` and proxies API requests to `http://localhost:8125`.

### Build for Production

```bash
npm run build
```

Build output goes to `goli/web/` directory, which is served by the Go backend.

## Project Structure

```
frontend/
├── src/
│   ├── api/
│   │   └── client.js          # API client & WebSocket
│   ├── components/
│   │   ├── Dashboard.vue      # Main dashboard
│   │   ├── Jobs.vue           # Jobs list
│   │   ├── Pipelines.vue      # Pipelines list
│   │   ├── Settings.vue       # User & system settings
│   │   ├── Login.vue          # Authentication
│   │   ├── SetupWizard.vue    # Initial setup
│   │   ├── LogsView.vue       # Logs viewer
│   │   ├── Modal.vue          # Reusable modal
│   │   ├── TextInput.vue      # Form input component
│   │   ├── FormField.vue      # Form field wrapper
│   │   ├── Alert.vue          # Alert component
│   │   ├── StatusBadge.vue    # Status badge
│   │   └── ToggleSwitch.vue   # Toggle switch
│   ├── App.vue                # Main app component
│   ├── main.js                # Entry point
│   └── style.css              # Tailwind imports
├── package.json
└── vite.config.js
```

## UI Components

### Reusable Components

- **Modal**: Reusable modal dialog with animations
- **TextInput**: Styled input with error handling
- **FormField**: Form field wrapper with labels
- **Alert**: Success/error/warning alerts
- **StatusBadge**: Status indicators for jobs/steps
- **ToggleSwitch**: Toggle switch for settings

### Pages

- **Dashboard**: Overview with stats and recent jobs
- **Jobs**: Job list with filtering and actions
- **Pipelines**: Pipeline management and execution
- **Settings**: User profile, 2FA, and system config
- **Login**: Authentication with 2FA support

## API Integration

All API calls are handled through `src/api/client.js`:

```javascript
import { getJobs, createJob, getPipelines } from './api/client'
```

### WebSocket

Real-time updates via WebSocket:

```javascript
const ws = new WebSocket('ws://localhost:8125/ws')
ws.onmessage = (event) => {
  const message = JSON.parse(event.data)
  // Handle updates
}
```

## Styling

Uses Tailwind CSS with custom theme:

- Primary colors: Blue scale
- Components: Cards, buttons, forms
- Responsive: Mobile-first design
- Dark terminal: Logs viewer with dark theme

## Deployment

1. Build the frontend: `npm run build`
2. Output is automatically copied to `goli/web/`
3. Go backend serves static files from `./web/`
4. SPA routing handled by backend

## Features in Detail

### Dashboard
- Real-time job statistics
- Recent jobs overview
- Quick actions

### Jobs Management
- List all jobs with status
- View job details
- Step-by-step logs
- Cancel running jobs

### Pipeline Management
- Upload YAML pipelines
- View pipeline definitions
- Run pipelines
- Monitor execution

### Settings
- User profile management
- 2FA configuration (Email/SMS)
- System configuration
- SMTP settings

## Development Tips

- Hot module replacement (HMR) enabled
- API calls proxied automatically
- WebSocket connections handled automatically
- Check browser console for debugging

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)
