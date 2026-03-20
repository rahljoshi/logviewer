# LogViewer

LogViewer is a local Graylog-style viewer for uploaded log files. It parses mixed log formats through a Go backend and presents searchable, filterable results in a React frontend.

## Prerequisites

- Go 1.21+
- Node 18+

## Running Locally

```bash
# Terminal 1
cd backend
go run .

# Terminal 2
cd frontend
npm install
npm run dev
```

The frontend runs on `http://localhost:5173` and proxies `/api/*` to `http://localhost:8080`.

## Supported Log Formats

- JSON: `{"timestamp":"2026-03-21T01:00:00Z","level":"info","message":"user login successful"}`
- logfmt: `time=2026-03-21T01:00:00Z level=warn msg="database query slow" request_id=abc123`
- plaintext: `2026-03-21 01:00:00 ERROR request timed out after 3s`

## Project Structure

```text
logviewer/
├── backend/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── parser/
│   ├── store/
│   └── testdata/
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── hooks/
│   │   └── utils/
│   ├── index.html
│   └── vite.config.js
├── AGENTS.md
└── workflow/
```

## What's Coming Next

- Live log streaming via SSE
- Persistent storage behind the existing store abstraction
