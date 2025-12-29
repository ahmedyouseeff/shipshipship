# Docker Database Status Issue - Fix

## Problem
When running the application in Docker, default statuses like "Feedback", "Proposed", "Backlog", etc. were being removed because the theme initialization was creating statuses from theme categories instead of preserving the default ones.

## Solution
The code has been updated to:

1. **Seed Default Statuses First**: A new function `SeedDefaultStatuses()` has been added that creates default statuses (Backlog, Proposed, Feedback, In Progress, Released, Archived) when the database is empty.

2. **Preserve Existing Statuses**: The theme initialization now checks if statuses already exist before creating new ones from theme categories.

## How It Works

1. When the database is initialized, `SeedDefaultStatuses()` is called first
2. This creates default statuses if the database is empty
3. When a theme is applied, `CreateDefaultStatusesFromTheme()` checks if statuses exist
4. If statuses exist, it skips creating new ones and only creates mappings

## To Fix Your Existing Database

If you already have a Docker volume with a database that's missing statuses, you have two options:

### Option 1: Reset the Database (Recommended for Development)

```bash
# Stop the container
docker-compose down

# Remove the volume (WARNING: This deletes all data)
docker volume rm shipshipship_shipshipship_data

# Start fresh
docker-compose up -d
```

### Option 2: Manually Add Statuses via API

1. Access your running container:
```bash
docker exec -it <container_name> /bin/bash
```

2. Or use the admin panel at `http://localhost:8088/admin` to create the missing statuses:
   - Go to Events page
   - Create statuses: "Backlog", "Proposed", "Feedback", "In Progress", "Released", "Archived"

### Option 3: Use SQLite to Add Statuses

```bash
# Access the database
docker exec -it <container_name> sqlite3 /app/data/changelog.db

# Then run SQL commands to insert statuses
# (This is complex and not recommended - use Option 1 or 2 instead)
```

## Building and Running

1. **Build the image:**
```bash
docker-compose build
```

2. **Start the container:**
```bash
docker-compose up -d
```

3. **Check logs:**
```bash
docker-compose logs -f shipshipship
```

You should see messages like:
```
Created default status: Backlog (order: 0)
Created default status: Proposed (order: 1)
Created default status: Feedback (order: 2)
...
```

## Important Notes

- The database is persisted in a Docker volume (`shipshipship_data`)
- Default statuses are only created if the database is empty
- Once statuses exist, they won't be overwritten by theme initialization
- Make sure to backup your data before resetting the database

## Verification

After starting the container, check that statuses exist:

1. Log into the admin panel: `http://localhost:8088/admin`
2. Go to the Events page
3. You should see the default statuses: Backlog, Proposed, Feedback, In Progress, Released, Archived


