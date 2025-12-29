#!/bin/bash

# Script to build a theme and apply it to ShipShipShip
# Usage: ./build-and-apply-theme.sh /path/to/theme/repo

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if theme path is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Theme repository path is required${NC}"
    echo "Usage: ./build-and-apply-theme.sh /path/to/theme/repo"
    exit 1
fi

THEME_REPO="$1"
THEME_BUILD_DIR="$THEME_REPO/build"
TARGET_DIR="./data/themes/current"
BACKUP_DIR="./data/themes/backup"

# Check if theme repository exists
if [ ! -d "$THEME_REPO" ]; then
    echo -e "${RED}Error: Theme repository not found at: $THEME_REPO${NC}"
    exit 1
fi

echo -e "${GREEN}Building and applying theme from: $THEME_REPO${NC}"

# Step 1: Build the theme
echo -e "${YELLOW}Step 1: Building theme...${NC}"
cd "$THEME_REPO"

# Check if package.json exists
if [ ! -f "package.json" ]; then
    echo -e "${RED}Error: package.json not found. Is this a valid theme repository?${NC}"
    exit 1
fi

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}Installing dependencies...${NC}"
    npm install
fi

# Build the theme
echo -e "${YELLOW}Running npm run build...${NC}"
npm run build

# Check if build was successful
if [ ! -d "$THEME_BUILD_DIR" ]; then
    echo -e "${RED}Error: Build directory not found. Build may have failed.${NC}"
    exit 1
fi

if [ ! -f "$THEME_BUILD_DIR/index.html" ]; then
    echo -e "${RED}Error: index.html not found in build directory.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Theme built successfully${NC}"

# Step 2: Create backup
echo -e "${YELLOW}Step 2: Creating backup of current theme...${NC}"
cd - > /dev/null  # Return to ShipShipShip root

if [ -d "$TARGET_DIR" ]; then
    mkdir -p "$BACKUP_DIR"
    rm -rf "$BACKUP_DIR"
    cp -r "$TARGET_DIR" "$BACKUP_DIR"
    echo -e "${GREEN}✓ Backup created at: $BACKUP_DIR${NC}"
else
    echo -e "${YELLOW}No existing theme to backup${NC}"
fi

# Step 3: Apply the new theme
echo -e "${YELLOW}Step 3: Applying new theme...${NC}"

# Remove old theme
rm -rf "$TARGET_DIR"
mkdir -p "$TARGET_DIR"

# Copy build files
cp -r "$THEME_BUILD_DIR"/* "$TARGET_DIR/"

# Copy theme.json if it exists
if [ -f "$THEME_REPO/theme.json" ]; then
    cp "$THEME_REPO/theme.json" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Copied theme.json${NC}"
fi

# Copy theme-manifest.json if it exists
if [ -f "$THEME_REPO/theme-manifest.json" ]; then
    cp "$THEME_REPO/theme-manifest.json" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Copied theme-manifest.json${NC}"
else
    # Create a basic manifest if it doesn't exist
    cat > "$TARGET_DIR/theme-manifest.json" << EOF
{
  "name": "custom-theme",
  "version": "v1.0.0",
  "built_with": "SvelteKit",
  "build_date": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "api_compatibility": "v1",
  "repository": "local-development"
}
EOF
    echo -e "${YELLOW}✓ Created basic theme-manifest.json${NC}"
fi

echo -e "${GREEN}✓ Theme applied successfully!${NC}"
echo ""
echo -e "${GREEN}Your theme is now available at: http://localhost:8080${NC}"
echo -e "${YELLOW}Note: You may need to restart the ShipShipShip backend for changes to take effect${NC}"
echo ""
echo -e "Backup location: $BACKUP_DIR"

