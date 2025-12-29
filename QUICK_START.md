# 🚀 Quick Start Guide - ShipShipShip

Simple step-by-step instructions to run the app.

## Prerequisites

- **Go 1.21+** installed
- **Node.js 20+** and npm installed
- **Git** (if cloning from repository)
- **Docker** (optional, only if using Docker method)

---

## ⚠️ Docker Not Running?

If you see: `Cannot connect to the Docker daemon... Is the docker daemon running?`

**Quick Fix:**
1. **Start Docker Desktop** (if installed)
   - On Linux: `sudo systemctl start docker` or start Docker Desktop app
   - On Windows/Mac: Open Docker Desktop application

2. **Or use Local Development instead** (no Docker needed) - see Option 1 below

---

## Option 1: Local Development (No Docker Required) ⭐ Recommended if Docker isn't running

### Step 1: Install Backend Dependencies
```bash
cd backend
go mod download
cd ..
```

### Step 2: Install Admin Dependencies
```bash
cd admin
npm install
cd ..
```

### Step 3: Run the App

**Option A: Use the development script (recommended)**
```bash
./start-dev.sh
```

This will:
- Build the backend
- Build the admin panel
- Start both servers

**Access URLs:**
- Backend API: http://localhost:8080
- Admin Panel: http://localhost:5173

**Option B: Quick start (backend only)**
```bash
./quick-start.sh
```

Then manually start the admin:
```bash
cd admin
npm run dev
```

### Step 4: Access the App
- **Admin Panel:** http://localhost:5173 (dev) or http://localhost:8080/admin (production)
- **Backend API:** http://localhost:8080

**Default login:**
- Username: `admin`
- Password: `admin`

**Done!** ✅

---

## Option 2: Docker (Requires Docker Running)

### Step 1: Run with Docker
```bash
docker run -d \
  -p 8080:8080 \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=changeme \
  -e JWT_SECRET=your-secret-key \
  -v shipshipship_data:/app/data \
  nelkinsky/shipshipship:latest
```

### Step 2: Access the App
Open your browser and go to: **http://localhost:8080/admin**

**Default login:**
- Username: `admin`
- Password: `changeme`

**Done!** ✅

---

## Option 3: Docker Compose (Requires Docker Running)

### Step 1: Create docker-compose.yml
Create a file named `docker-compose.yml` with:
```yaml
version: "3.8"
services:
  shipshipship:
    image: nelkinsky/shipshipship:latest
    ports:
      - "8080:8080"
    environment:
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=changeme
      - JWT_SECRET=your-secret-key
    volumes:
      - shipshipship_data:/app/data
    restart: unless-stopped

volumes:
  shipshipship_data:
```

### Step 2: Start the App
```bash
docker-compose up -d
```

### Step 3: Access the App
Open your browser and go to: **http://localhost:8080/admin**

**Done!** ✅

---

## Option 4: Manual Setup (Step by Step)

### Step 1: Build Backend
```bash
cd backend
go build -o main .
cd ..
```

### Step 2: Build Admin
```bash
cd admin
npm install
npm run build
cd ..
```

### Step 3: Start Backend
```bash
cd backend
./main
```

The backend will run on **http://localhost:8080**

### Step 4: Start Admin (in a new terminal)
```bash
cd admin
npm run dev
```

The admin will run on **http://localhost:5173**

**Done!** ✅

---

## Environment Variables

You can customize the app with these environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `ADMIN_USERNAME` | `admin` | Admin username |
| `ADMIN_PASSWORD` | `admin` | Admin password |
| `JWT_SECRET` | `your-secret-key-change-in-production` | JWT signing key |
| `BASE_URL` | _(auto-detected)_ | Base URL for your instance |
| `PORT` | `8080` | Server port |
| `GIN_MODE` | `debug` | `debug` or `release` |
| `DB_PATH` | `./data/changelog.db` | Database path |

**Example:**
```bash
export ADMIN_USERNAME=myadmin
export ADMIN_PASSWORD=mypassword
export JWT_SECRET=my-secret-key
./start-dev.sh
```

---

## Troubleshooting

### Docker daemon not running
**Error:** `Cannot connect to the Docker daemon at unix://.../docker.sock. Is the docker daemon running?`

**Solutions:**
1. **Start Docker Desktop** application
2. **Or start Docker service** (Linux):
   ```bash
   sudo systemctl start docker
   sudo systemctl enable docker  # Enable on boot
   ```
3. **Or use Local Development** (Option 1) - no Docker needed!

### Port 8080 already in use
```bash
# Kill the process using port 8080
lsof -ti:8080 | xargs kill -9
```

### Backend won't start
- Check `backend.log` for errors
- Make sure Go 1.21+ is installed: `go version`
- Make sure dependencies are installed: `cd backend && go mod download`

### Admin won't start
- Check `admin.log` for errors
- Make sure Node.js 20+ is installed: `node --version`
- Reinstall dependencies: `cd admin && rm -rf node_modules && npm install`

### Database issues
- The database is created automatically at `./data/changelog.db`
- Make sure the `data` directory exists and is writable

---

## Next Steps

1. **Login** to the admin panel
2. **Create events** in the Kanban board
3. **Install a theme** at `/admin/appearance/theme`
4. **Configure newsletter** at `/admin/newsletter/settings`
5. **Map statuses** to theme categories

---

## Useful Commands

```bash
# View backend logs
tail -f backend.log

# View admin logs
tail -f admin.log

# Rebuild everything
./start-dev.sh --rebuild

# Stop all servers (if using start-dev.sh)
# Press Ctrl+C in the terminal
```

---

**Need help?** Check the main [README.md](README.md) for more details.

