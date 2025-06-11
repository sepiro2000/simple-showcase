#!/bin/bash

# Exit on error
set -e

# Disable history expansion
set +H

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Show help message
show_help() {
    echo "Backend Deployment Script"
    echo "Usage: ./deploy-backend.sh [ENV_VAR=VALUE]..."
    echo ""
    echo "Example: ./deploy-backend.sh 'DB_PASSWORD=abc123' 'WRITE_DB_HOST=localhost'"
    echo ""
    echo "Available environment variables:"
    echo "  WRITE_DB_HOST    - Write database host (default: 127.0.0.1)"
    echo "  READ_DB_HOST     - Read database host (default: 127.0.0.1)"
    echo "  DB_PORT         - Database port (default: 3306)"
    echo "  DB_USER         - Database user (default: showcase_user)"
    echo "  DB_PASSWORD     - Database password (default: YOUR_APP_PASSWORD_HERE)"
    echo "  DB_NAME         - Database name (default: simple_showcase)"
    echo "  APP_PORT        - Application port (default: 8080)"
    echo ""
    echo "Note: Use single quotes around values containing special characters"
    echo "Example: ./deploy-backend.sh 'DB_PASSWORD=abc!@#'"
    exit 0
}

# Check for help argument
if [ "$1" = "help" ] || [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    show_help
fi

# Default values
DEFAULT_WRITE_DB_HOST="127.0.0.1"
DEFAULT_READ_DB_HOST="127.0.0.1"
DEFAULT_DB_PORT="3306"
DEFAULT_DB_USER="showcase_user"
DEFAULT_DB_PASSWORD="YOUR_APP_PASSWORD_HERE"
DEFAULT_DB_NAME="simple_showcase"
DEFAULT_APP_PORT="8080"

# Get current user and home directory
CURRENT_USER=$(whoami)
HOME_DIR=$(eval echo ~$CURRENT_USER)
APP_PATH="/home/$CURRENT_USER/app/simple-showcase/backend"
DEPLOY_PATH="/opt/simple-showcase-backend"

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

# Process environment variables from arguments
for arg in "$@"; do
    if [[ $arg == *"="* ]]; then
        key="${arg%%=*}"
        value="${arg#*=}"
        case $key in
            WRITE_DB_HOST) DEFAULT_WRITE_DB_HOST="$value" ;;
            READ_DB_HOST) DEFAULT_READ_DB_HOST="$value" ;;
            DB_PORT) DEFAULT_DB_PORT="$value" ;;
            DB_USER) DEFAULT_DB_USER="$value" ;;
            DB_PASSWORD) DEFAULT_DB_PASSWORD="$value" ;;
            DB_NAME) DEFAULT_DB_NAME="$value" ;;
            APP_PORT) DEFAULT_APP_PORT="$value" ;;
            *) echo -e "${RED}Warning: Unknown environment variable $key${NC}" ;;
        esac
    fi
done

echo "Starting backend deployment..."
echo -e "Using environment variables:"
echo -e "  WRITE_DB_HOST: ${GREEN}$DEFAULT_WRITE_DB_HOST${NC}"
echo -e "  READ_DB_HOST: ${GREEN}$DEFAULT_READ_DB_HOST${NC}"
echo -e "  DB_PORT: ${GREEN}$DEFAULT_DB_PORT${NC}"
echo -e "  DB_USER: ${GREEN}$DEFAULT_DB_USER${NC}"
echo -e "  DB_PASSWORD: ${GREEN}****${NC}"
echo -e "  DB_NAME: ${GREEN}$DEFAULT_DB_NAME${NC}"
echo -e "  APP_PORT: ${GREEN}$DEFAULT_APP_PORT${NC}"

# Install required packages
echo -n "Installing required packages... "
sudo dnf update -y && sudo dnf install -y golang
check_status "Failed to install required packages"

# Navigate to backend directory
echo -n "Navigating to backend directory... "
cd "$APP_PATH"
check_status "Failed to navigate to backend directory"

# Build the application
echo -n "Building the application... "
go build -o simple-showcase-backend ./
check_status "Failed to build the application"

# Create deployment directory
echo -n "Creating deployment directory... "
sudo mkdir -p "$DEPLOY_PATH"
check_status "Failed to create deployment directory"

# Move binary to deployment directory
echo -n "Moving binary to deployment directory... "
sudo mv simple-showcase-backend "$DEPLOY_PATH/"
check_status "Failed to move binary"

# Update service file with environment variables
echo -n "Updating service file... "
SERVICE_FILE="$APP_PATH/deploy/simple-showcase-backend.service"
TEMP_FILE=$(mktemp)

# Escape special characters in values
escape_value() {
    echo "$1" | sed 's/[\/&]/\\&/g'
}

WRITE_DB_HOST_ESC=$(escape_value "$DEFAULT_WRITE_DB_HOST")
READ_DB_HOST_ESC=$(escape_value "$DEFAULT_READ_DB_HOST")
DB_PORT_ESC=$(escape_value "$DEFAULT_DB_PORT")
DB_USER_ESC=$(escape_value "$DEFAULT_DB_USER")
DB_PASSWORD_ESC=$(escape_value "$DEFAULT_DB_PASSWORD")
DB_NAME_ESC=$(escape_value "$DEFAULT_DB_NAME")
APP_PORT_ESC=$(escape_value "$DEFAULT_APP_PORT")

sed "s|Environment=\"WRITE_DB_HOST=.*\"|Environment=\"WRITE_DB_HOST=$WRITE_DB_HOST_ESC\"|" "$SERVICE_FILE" > "$TEMP_FILE"
sed -i "s|Environment=\"READ_DB_HOST=.*\"|Environment=\"READ_DB_HOST=$READ_DB_HOST_ESC\"|" "$TEMP_FILE"
sed -i "s|Environment=\"DB_PORT=.*\"|Environment=\"DB_PORT=$DB_PORT_ESC\"|" "$TEMP_FILE"
sed -i "s|Environment=\"DB_USER=.*\"|Environment=\"DB_USER=$DB_USER_ESC\"|" "$TEMP_FILE"
sed -i "s|Environment=\"DB_PASSWORD=.*\"|Environment=\"DB_PASSWORD=$DB_PASSWORD_ESC\"|" "$TEMP_FILE"
sed -i "s|Environment=\"DB_NAME=.*\"|Environment=\"DB_NAME=$DB_NAME_ESC\"|" "$TEMP_FILE"
sed -i "s|Environment=\"APP_PORT=.*\"|Environment=\"APP_PORT=$APP_PORT_ESC\"|" "$TEMP_FILE"

sudo cp "$TEMP_FILE" /etc/systemd/system/simple-showcase-backend.service
rm "$TEMP_FILE"
check_status "Failed to update service file"

# Reload systemd
echo -n "Reloading systemd... "
sudo systemctl daemon-reload
check_status "Failed to reload systemd"

# Restart service
echo -n "Restarting service... "
sudo systemctl restart simple-showcase-backend
check_status "Failed to restart service"

# Enable service
echo -n "Enabling service... "
sudo systemctl enable simple-showcase-backend
check_status "Failed to enable service"

echo -e "\n${GREEN}✓ Backend deployment completed successfully!${NC}"
