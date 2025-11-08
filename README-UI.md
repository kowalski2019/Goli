# Goli UI Development Guide

## Current UI Structure

The current UI is a simple HTML/JavaScript dashboard located in `goli/web/index.html`. This serves as a temporary interface while the full Vue.js application is being developed.

## Future UI Architecture

### Technology Stack
- **ViteJS** - Build tool and dev server
- **Vue 3** - Frontend framework
- **JavaScript** - Programming language
- **Tailwind CSS** - Utility-first CSS framework

### Directory Structure

```
goli/
├── web/                    # Static files served by Go server
│   ├── index.html          # Current simple UI (temporary)
│   └── example-pipeline.yaml
│
frontend/                   # Future Vue.js application (to be created)
├── src/
│   ├── components/
│   ├── views/
│   ├── router/
│   ├── store/
│   └── main.js
├── public/
├── package.json
├── vite.config.js
└── tailwind.config.js
```

### Build Process

1. **Development**: Run ViteJS dev server for hot-reload development
2. **Production**: Build with `npm run build` which outputs to `frontend/dist/`
3. **Deployment**: Copy `frontend/dist/*` contents to `goli/web/` directory

### Integration Steps

When ready to integrate the Vue.js UI:

1. Create the Vue.js project in a `frontend/` directory
2. Configure ViteJS to output to `goli/web/` or copy build output
3. Update `goli/main.go` to serve static files from `./web/` (already configured)
4. Ensure API endpoints remain at `/api/v1/*` (already configured)
5. WebSocket endpoint is at `/ws` (already configured)

### API Endpoints

All API endpoints are prefixed with `/api/v1/`:

- `GET /api/v1/jobs` - List jobs
- `POST /api/v1/jobs` - Create job
- `GET /api/v1/jobs/{id}` - Get job details
- `GET /api/v1/pipelines` - List pipelines
- `POST /api/v1/pipelines` - Create pipeline (JSON)
- `POST /api/v1/pipelines/upload` - Upload pipeline (YAML file)
- `GET /api/v1/pipelines/{id}` - Get pipeline details
- `POST /api/v1/pipelines/{id}/run` - Run pipeline

### WebSocket

- Endpoint: `ws://localhost:8125/ws` (or `wss://` for HTTPS)
- Message types:
  - `job_update` - Job status changed
  - `stats_update` - Statistics updated

### Current Features

✅ Job management (create, list, view)
✅ Pipeline upload via YAML file
✅ Pipeline execution
✅ Real-time updates via WebSocket
✅ Dashboard with statistics

### Future Enhancements

- Full Vue.js SPA with routing
- Pipeline editor (visual + YAML)
- Job logs viewer
- Pipeline history and rollback
- User authentication UI
- Settings and configuration UI

