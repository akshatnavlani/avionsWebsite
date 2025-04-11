#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${YELLOW}$1...${NC}"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error and exit
print_error() {
    echo -e "${RED}✗ $1${NC}"
    exit 1
}

# Check if .env exists
print_status "Checking environment configuration"
if [ ! -f .env ]; then
    print_error "Missing .env file. Please copy .env.example to .env and configure it"
fi
print_success "Environment configuration found"

# Set production mode
print_status "Setting production mode"
sed -i 's/ENV=development/ENV=production/' .env
print_success "Production mode set"

# Clean and get dependencies
print_status "Installing dependencies"
go mod tidy
if [ $? -ne 0 ]; then
    print_error "Failed to install dependencies"
fi
print_success "Dependencies installed"

# Build the application
print_status "Building application"
go build -o server
if [ $? -ne 0 ]; then
    print_error "Failed to build application"
fi
print_success "Application built successfully"

# Run database migrations
print_status "Running database migrations"
./server migrate
if [ $? -ne 0 ]; then
    print_error "Failed to run database migrations"
fi
print_success "Database migrations completed"

print_success "Deployment completed successfully!"
echo -e "\nYou can now run the server with: ${GREEN}./server${NC}" 