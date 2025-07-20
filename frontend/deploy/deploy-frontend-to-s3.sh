#!/bin/bash

# Exit on error
set -e

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Show help message
show_help() {
    echo "Frontend S3 Deployment Script"
    echo "Usage: ./deploy-frontend-to-s3.sh <S3_BUCKET_NAME> [API_URL] [NODE_VERSION]"
    echo "Example: ./deploy-frontend-to-s3.sh my-bucket-name https://api.example.com 20.11.1"
    echo ""
    echo "Arguments:"
    echo "  S3_BUCKET_NAME - (required) Target S3 bucket name for deployment"
    echo "  API_URL        - Backend API URL (default: http://localhost:8080)"
    echo "  NODE_VERSION   - Node.js version to install (default: 20.11.1)"
    exit 0
}

# Check for help argument
if [ "$1" = "help" ] || [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    show_help
fi

# Check for required S3 bucket argument
if [ -z "$1" ]; then
    echo -e "${RED}Error: S3_BUCKET_NAME is required.${NC}"
    show_help
fi

# Default values
DEFAULT_API_URL="http://localhost:8080"
DEFAULT_NODE_VERSION="20.11.1"  # LTS version
DEFAULT_APP_PATH="/home/ec2-user/app/simple-showcase/frontend"

# Get current user and home directory
CURRENT_USER=$(whoami)
HOME_DIR=$(eval echo ~$CURRENT_USER)

# Parse command line arguments
S3_BUCKET="$1"
API_URL=${2:-$DEFAULT_API_URL}
NODE_VERSION=${3:-$DEFAULT_NODE_VERSION}
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

echo "Starting frontend S3 deployment..."
echo -e "Using S3 bucket: ${GREEN}$S3_BUCKET${NC}"
echo -e "Using API URL: ${GREEN}$API_URL${NC}"
echo -e "Using Node.js version: ${GREEN}$NODE_VERSION${NC}"
echo -e "App path: ${GREEN}$APP_PATH${NC}"

# Check for AWS CLI
if ! command -v aws &> /dev/null; then
    echo -n "Installing AWS CLI... "
    sudo dnf install -y awscli
    check_status "Failed to install AWS CLI"
fi

# Check AWS credentials
if ! aws sts get-caller-identity &> /dev/null; then
    handle_error "AWS CLI is not configured. Please run 'aws configure' first."
fi

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

# Install Node.js
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
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' "s|VITE_API_BASE_URL=.*|VITE_API_BASE_URL=$API_URL|" .env
else
  sed -i "s|VITE_API_BASE_URL=.*|VITE_API_BASE_URL=$API_URL|" .env
fi
check_status "Failed to configure API URL"

# Install dependencies
echo -n "Installing dependencies... "
npm install
check_status "Failed to install dependencies"

# Build the application
echo -n "Building the application... "
npm run build
check_status "Failed to build the application"

# Upload to S3
echo -n "Uploading build output to S3... "
aws s3 sync ./dist/ s3://$S3_BUCKET/ --delete
check_status "Failed to upload to S3"

echo -e "\n${GREEN}✓ Frontend S3 deployment completed successfully!${NC}"