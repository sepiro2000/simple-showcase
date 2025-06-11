#!/bin/bash

# Exit on error
set -e

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default values
DEFAULT_API_URL="http://localhost:8080"
DEFAULT_NODE_VERSION="20.11.1"  # LTS version
DEFAULT_DEPLOY_PATH="/usr/share/nginx/html/simple-showcase-frontend"
DEFAULT_APP_PATH="/home/ec2-user/app/simple-showcase/frontend"

# Get current user and home directory
CURRENT_USER=$(whoami)
HOME_DIR=$(eval echo ~$CURRENT_USER)

# Parse command line arguments
API_URL=${1:-$DEFAULT_API_URL}
NODE_VERSION=${2:-$DEFAULT_NODE_VERSION}
DEPLOY_PATH=${3:-$DEFAULT_DEPLOY_PATH}
APP_PATH=${4:-$DEFAULT_APP_PATH}

# Function to handle errors
handle_error() {
    echo -e "${RED}Error: $1${NC}"
    exit 1
}

# Function to check command status
check_status() {
    if [ $? -ne 0 ]; then
        handle_error "$1"
    else
        echo -e "${GREEN}✓ PASS${NC}"
    fi
}

echo "Starting frontend deployment..."
echo -e "Using API URL: ${GREEN}$API_URL${NC}"
echo -e "Using Node.js version: ${GREEN}$NODE_VERSION${NC}"
echo -e "Deploy path: ${GREEN}$DEPLOY_PATH${NC}"
echo -e "App path: ${GREEN}$APP_PATH${NC}"

# Install required packages
echo -n "Installing required packages... "
sudo dnf update -y && sudo dnf install -y curl git
check_status "Failed to install required packages"

# Install NVM
echo -n "Installing NVM... "
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
check_status "Failed to install NVM"

# Source NVM (try multiple possible locations)
echo -n "Sourcing NVM... "
if [ -f "$HOME_DIR/.nvm/nvm.sh" ]; then
    source "$HOME_DIR/.nvm/nvm.sh"
elif [ -f "$HOME_DIR/.bashrc" ]; then
    source "$HOME_DIR/.bashrc"
else
    handle_error "Could not find NVM installation"
fi
check_status "Failed to source NVM"

# Install specific Node.js version
echo -n "Installing Node.js $NODE_VERSION... "
nvm install $NODE_VERSION
nvm use $NODE_VERSION
check_status "Failed to install Node.js"

# Verify Node.js and npm installation
echo -n "Verifying Node.js and npm installation... "
node -v
npm -v
check_status "Failed to verify Node.js and npm installation"

# Navigate to frontend directory
echo -n "Navigating to frontend directory... "
cd "$APP_PATH"
check_status "Failed to navigate to frontend directory"

# Setup environment
echo -n "Setting up environment... "
cp .env.example .env
check_status "Failed to copy environment file"

# Update API URL in .env file
echo -n "Configuring API URL... "
sed -i "s|VITE_API_BASE_URL=.*|VITE_API_BASE_URL=$API_URL|" .env
check_status "Failed to configure API URL"

# Install dependencies
echo -n "Installing dependencies... "
npm install
check_status "Failed to install dependencies"

# Build the application
echo -n "Building the application... "
npm run build
check_status "Failed to build the application"

# Create nginx directory
echo -n "Creating nginx directory... "
sudo mkdir -p "$DEPLOY_PATH"
check_status "Failed to create nginx directory"

# Copy build files
echo -n "Copying build files... "
sudo cp -a "$APP_PATH/dist/." "$DEPLOY_PATH/"
check_status "Failed to copy build files"

# Set permissions
echo -n "Setting permissions... "
sudo chown -R nginx:nginx "$DEPLOY_PATH"
check_status "Failed to set permissions"

# Copy nginx configuration
echo -n "Copying nginx configuration... "
sudo cp -a "$APP_PATH/deploy/simple-showcase-frontend.conf" /etc/nginx/conf.d/simple-showcase-frontend.conf
check_status "Failed to copy nginx configuration"

# Test nginx configuration
echo -n "Testing nginx configuration... "
sudo nginx -t
check_status "Failed to test nginx configuration"

# Restart nginx
echo -n "Restarting nginx... "
sudo systemctl restart nginx
check_status "Failed to restart nginx"

echo -e "\n${GREEN}✓ Frontend deployment completed successfully!${NC}"
echo -e "Usage: ./deploy-frontend.sh [API_URL] [NODE_VERSION] [DEPLOY_PATH] [APP_PATH]"
echo -e "Example: ./deploy-frontend.sh \"https://api.example.com\" \"20.11.1\" \"/usr/share/nginx/html/myapp\" \"/home/ec2-user/app/myapp/frontend\""