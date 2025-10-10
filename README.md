<p align="center">
<img width="2000" height="500" alt="nabd second" src="https://github.com/user-attachments/assets/158bcfcc-5871-4252-afe5-ab9ea661840a" /> 

<p align="center">Lightweight Open-Source Container Observability & Auto-Healing Tool</p>

<p align="center">Nabd is a self-contained observability and auto-healing toolkit for Docker containers that combines metrics collection, log monitoring, alerting, and container auto-healing in a single, easy-to-deploy package.</p>

<div align="center">
  <img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License" />
  <img src="https://img.shields.io/badge/version-v0.1.0-green.svg" alt="Version" />
  <img src="https://img.shields.io/badge/Go-1.21+-blue.svg" alt="Go" />
  <img src="https://img.shields.io/badge/React-18+-blue.svg" alt="React" />
</div>


## Features

### Metrics Collection
- Real-time monitoring of CPU, Memory, Network, and Disk usage
- Historical metrics storage in SQLite database
- REST API endpoints for metrics data
- Automated data collection every 15 seconds

### Log Monitoring
- Container log collection and viewing
- Real-time log streaming
- Configurable log history (50-500 lines)
- Beautiful terminal-style log viewer

### Auto-Healing
- Automatic detection of stopped/unhealthy containers
- Smart container restart capabilities
- Configurable restart policies and limits
- Comprehensive event logging
- Manual trigger support

### Intelligent Alerting
- CPU and memory threshold alerts
- Container state change notifications
- Customizable alert thresholds via config
- Visual alert dashboard

### Modern Dashboard
- Beautiful dark-mode UI built with React + TailwindCSS
- Real-time container status overview
- Interactive metrics visualization
- Mobile-responsive design
- Token-based authentication

### Security
- JWT-based authentication
- Configurable admin token
- Secure API endpoints
- Docker socket access control

## Getting Started

### One-Command Deployment

```bash
docker run -d \
  --name nabd \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v nabd_data:/data \
  -e NABD_ADMIN_TOKEN=your-secure-token \
  nabd/nabd:v0.1.0
```


## Installation Options

### Option 1: Docker Compose (Full Stack)
```bash
git clone https://github.com/your-username/nabd.git
cd nabd
docker-compose up -d
```

### Option 2: Docker Run (Single Container)
```bash
docker run -d \
  --name nabd \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v nabd_data:/data \
  nabd/nabd:latest
```

### Option 3: Build from Source
```bash
# Backend
cd backend
go mod download
go build -o nabd main.go

# Frontend
cd ../frontend
npm install
npm run build

# Run
./backend/nabd
```

## Configuration

### Environment Variables
```bash
NABD_ADMIN_TOKEN=your-secure-token    # Dashboard authentication token
NABD_DB_PATH=/data/nabd.db           # SQLite database path
DOCKER_HOST=unix:///var/run/docker.sock # Docker socket path
```

### Configuration File (`config.yaml`)
```yaml
database:
  path: "./nabd.db"

docker:
  host: "unix:///var/run/docker.sock"

auth:
  admin_token: "nabd-admin-token"

alerts:
  cpu_threshold: 90.0      # CPU % threshold for alerts
  memory_threshold: 90.0   # Memory % threshold for alerts
  restart_limit: 3         # Max restarts before alerting
```

## API Reference

### Authentication
```bash
POST /api/auth/login
{
  "token": "your-admin-token"
}
```

### Container Metrics
```bash
GET /api/containers          # List all containers
GET /api/metrics             # Current metrics for all containers
GET /api/metrics/:id/history # Historical metrics for container
GET /api/logs?container=name # Container logs
POST /api/containers/:name/restart # Restart container
```

### Auto-Healing
```bash
GET /api/autoheal/history    # Auto-heal event history
POST /api/autoheal/trigger   # Manually trigger auto-heal check
```

### Alerts
```bash
GET /api/alerts              # Get active alerts
```

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   React App     │    │   Go Backend     │    │  Docker Engine  │
│                 │    │                  │    │                 │
│ • Dashboard     │◄──►│ • Gin Framework  │◄──►│ • Containers    │
│ • Log Viewer    │    │ • REST API       │    │ • Metrics       │
│ • Auto-heal UI  │    │ • SQLite DB      │    │ • Logs          │
│ • Alerts        │    │ • Auto-healing   │    │ • Events        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Tech Stack:**
- **Backend:** Go 1.21+ (Gin framework)
- **Frontend:** React 18 + TailwindCSS
- **Database:** SQLite (embedded)
- **Container:** Docker SDK for Go
- **Authentication:** JWT tokens


## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.


## Star History

If you find Nabd useful, please consider giving it a star on GitHub!

---

**Made with ❤️ by the Nabd Team**
