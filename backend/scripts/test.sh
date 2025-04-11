#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to make API calls
call_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4

    if [ -z "$headers" ]; then
        headers="Content-Type: application/json"
    fi

    curl -s -X "$method" "$BASE_URL$endpoint" -H "$headers" ${data:+-d "$data"}
}

# Function to print test results
print_result() {
    local test_name=$1
    local success=$2
    if [ "$success" = true ]; then
        echo -e "${GREEN}✓ $test_name${NC}"
    else
        echo -e "${RED}✗ $test_name${NC}"
    fi
}

echo "Starting API Tests..."

# Get admin token
echo -e "\n${BLUE}1. Testing Authentication${NC}"
ADMIN_PASSWORD=$(grep ADMIN_PASSWORD .env | cut -d '=' -f2 | tr -d '"' | tr -d '\r')
AUTH_RESPONSE=$(call_api "POST" "/api/auth/login" "{\"password\": \"$ADMIN_PASSWORD\"}")
TOKEN=$(echo "$AUTH_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get token${NC}"
    echo "Response: $AUTH_RESPONSE"
    exit 1
fi
print_result "Admin Authentication" true

# Test Member Operations
echo -e "\n${BLUE}2. Testing Member Operations${NC}"
MEMBER_RESPONSE=$(call_api "POST" "/api/members" '{
    "name": "Test Member",
    "position": "Software Engineer",
    "imageUrl": "https://example.com/image.jpg"
}' "Authorization: Bearer $TOKEN")

MEMBER_ID=$(echo "$MEMBER_RESPONSE" | grep -o '"id":"[^"]*"' | head -n 1 | cut -d'"' -f4)
[ ! -z "$MEMBER_ID" ] && print_result "Create Member" true || print_result "Create Member" false

# Test Blog Operations
echo -e "\n${BLUE}3. Testing Blog Operations${NC}"
BLOG_RESPONSE=$(call_api "POST" "/api/blogs" "{
    \"title\": \"Test Blog\",
    \"description\": \"Test Description\",
    \"markdownUrl\": \"https://example.com/blog.md\",
    \"authorId\": \"$MEMBER_ID\"
}" "Authorization: Bearer $TOKEN")

BLOG_ID=$(echo "$BLOG_RESPONSE" | grep -o '"id":"[^"]*"' | head -n 1 | cut -d'"' -f4)
[ ! -z "$BLOG_ID" ] && print_result "Create Blog" true || print_result "Create Blog" false

# Test Project Operations
echo -e "\n${BLUE}4. Testing Project Operations${NC}"
PROJECT_RESPONSE=$(call_api "POST" "/api/projects" '{
    "title": "Test Project",
    "description": "Test Description",
    "markdownUrl": "https://example.com/project.md",
    "imageUrl": "https://example.com/project.jpg"
}' "Authorization: Bearer $TOKEN")

PROJECT_ID=$(echo "$PROJECT_RESPONSE" | grep -o '"id":"[^"]*"' | head -n 1 | cut -d'"' -f4)
[ ! -z "$PROJECT_ID" ] && print_result "Create Project" true || print_result "Create Project" false

# Clean up
echo -e "\n${BLUE}5. Testing Cleanup Operations${NC}"

# Delete Blog
DELETE_BLOG=$(call_api "DELETE" "/api/blogs/$BLOG_ID" "" "Authorization: Bearer $TOKEN")
[ ! -z "$DELETE_BLOG" ] && print_result "Delete Blog" true || print_result "Delete Blog" false

# Delete Project
DELETE_PROJECT=$(call_api "DELETE" "/api/projects/$PROJECT_ID" "" "Authorization: Bearer $TOKEN")
[ ! -z "$DELETE_PROJECT" ] && print_result "Delete Project" true || print_result "Delete Project" false

# Delete Member
DELETE_MEMBER=$(call_api "DELETE" "/api/members/$MEMBER_ID" "" "Authorization: Bearer $TOKEN")
[ ! -z "$DELETE_MEMBER" ] && print_result "Delete Member" true || print_result "Delete Member" false

echo -e "\n${GREEN}Tests completed!${NC}" 