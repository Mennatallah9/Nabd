<p align="center">
<img width="2000" height="500" alt="nabd second" src="https://github.com/user-attachments/assets/158bcfcc-5871-4252-afe5-ab9ea661840a" />
</p>

<p align="center">Lightweight Open-Source Container Observability & Auto-Healing Tool</p>

<p align="center">Nabd is a self-contained observability and auto-healing toolkit for Docker containers that combines metrics collection, log monitoring, alerting, and container auto-healing in a single, easy-to-deploy package.</p>

<div align="center">
  <img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License" />
  <img src="https://img.shields.io/badge/version-v0.1.0-green.svg" alt="Version" />
  <img src="https://img.shields.io/badge/Go-1.21+-blue.svg" alt="Go" />
  <img src="https://img.shields.io/badge/React-18+-blue.svg" alt="React" />
</div>

> **💡 Fun Fact:** The word “Nabd” (نبض) in Arabic literally means “pulse” or “heartbeat”.


## Features

### Metrics Collection
- Real-time monitoring of CPU, Memory, Network, and Disk usage
- Historical metrics storage in SQLite database
- REST API endpoints for metrics data
- Automated data collection every 15 seconds

### Log Monitoring
- Container log collection and viewing
- Real-time log streaming
- **AI-powered log summarization using Chrome's built-in Gemini Nano**
- Configurable log history (50-500 lines)
- Beautiful terminal-style log viewer
- One-click intelligent log analysis

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

### Security
- JWT-based authentication
- Configurable admin token
- Secure API endpoints
- Docker socket access control

## Getting Started

### Installation Options

#### Option 1: One-Command Deployment

```bash
docker run -d \
  --name nabd \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v nabd_data:/data \
  -e NABD_ADMIN_TOKEN=your-secure-token \
  mennahaggag/nabd:latest
```
then navigate to `http://localhost:8080`

#### Option 2: Docker Compose Deployment

For a more customizable Docker deployment with persistent configuration:

1. Set up environment variables:
   ```bash
   cp .env.example .env
   ```
   
2. Edit the `.env` file with your settings

3. Copy and customize the configuration:
   ```bash
   cp config.yaml.example config.yaml
   ```

4. Deploy with Docker Compose:
   ```bash
   docker-compose up -d
   ```

This option provides:
- Environment-based configuration
- Persistent data volumes
- Example monitored container (nginx)
- Easy customization of settings


#### Option 3: Build from Source
1. Copy and customize the configuration:
    ```bash
    cp config.yaml.example backend/config.yaml
    ```
2. Build the application:

   **For Linux/macOS:**
   ```bash
   ./build.sh
   # then run
   ./backend/nabd
   ```
   
   **For Windows:**
   ```batch
   .\build.bat
   # then run
   .\backend\nabd.exe
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
GET /api/containers                    # List all containers
GET /api/metrics                       # Current metrics for all containers
GET /api/metrics/:id/history           # Historical metrics for container
GET /api/logs?container=name           # Container logs
POST /api/containers/:name/restart     # Restart container
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
- **AI:** Chrome's built-in Gemini Nano

## AI Log Summarization Setup

The AI log summarization feature uses Chrome's built-in Gemini Nano model, which runs locally in your browser for privacy and performance.

### Requirements:
- **Chrome 129+** (or Chrome Canary with experimental features enabled)
- **Supported OS:** Windows 10/11, macOS 13+, or Linux
- **Hardware:** At least 4GB VRAM (GPU) or 16GB RAM (CPU)
- **Storage:** At least 22GB free space for the model

### Setup:
1. **Enable the feature in Chrome:**
   - Open `chrome://flags/`
   - Search for "Gemini Nano" or "AI"
   - Enable the "Gemini Nano" flag
   - Restart Chrome

2. **Usage:**
   - Navigate to the Logs page in the dashboard
   - Select a container
   - Click the **"Summarize Logs"** button
   - Wait for the model to download (first time only)
   - View the AI-generated summary with key insights

### Features:
- **Completely Local:** All processing happens in your browser
- **Privacy:** No data sent to external servers
- **Key Insights:** Identifies errors, warnings, and patterns
- **Smart Analysis:** Provides actionable recommendations
- **No API Keys Required:** Uses Chrome's built-in AI

> **Note:** If you're not using Chrome, the summarization feature will show an appropriate message. All other Nabd features work normally in any modern browser.


## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.


## Star History

If you find Nabd useful, please consider giving it a star on GitHub!

---

**Made with ❤️ by the Nabd Team**
