#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080"

# Test files directory
TEST_FILES_DIR="test_files"
mkdir -p "$TEST_FILES_DIR"

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

# Function to test file upload
test_file_upload() {
    local file_path="$1"
    local file_type="$2"
    local expected_bucket="$3"

    if [ ! -f "$file_path" ]; then
        echo -e "${RED}Error: Test file not found: $file_path${NC}"
        return 1
    fi

    echo "Testing $file_type upload..."
    RESPONSE=$(curl -v -X POST "$BASE_URL/api/storage/upload" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: multipart/form-data" \
        -F "type=$file_type" \
        -F "file=@$file_path" 2>&1)

    # Extract URL and filename from response
    URL=$(echo "$RESPONSE" | grep -o '"url":"[^"]*"' | cut -d'"' -f4)
    FILENAME=$(echo "$RESPONSE" | grep -o '"filename":"[^"]*"' | cut -d'"' -f4)

    if [ -n "$URL" ] && [ -n "$FILENAME" ]; then
        echo -e "${GREEN}Upload successful!${NC}"
        echo "URL: $URL"
        echo "Filename: $FILENAME"
        return 0
    else
        echo -e "${RED}Upload failed${NC}"
        echo "Response: $RESPONSE"
        return 1
    fi
}

# Function to test database operations
test_database() {
    echo "Testing database operations..."

    # Test creating a user
    echo "Creating test user..."
    USER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/users" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser",
            "email": "test@example.com",
            "password": "testpass123"
        }')
    USER_ID=$(echo "$USER_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    print_result $? "Create user"

    # Test creating a post
    echo "Creating test post..."
    POST_RESPONSE=$(curl -s -X POST "$BASE_URL/api/posts" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"title\": \"Test Post\",
            \"content\": \"This is a test post content\",
            \"author_id\": \"$USER_ID\",
            \"status\": \"draft\"
        }")
    POST_ID=$(echo "$POST_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    print_result $? "Create post"

    # Test creating a comment
    echo "Creating test comment..."
    COMMENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/comments" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"content\": \"This is a test comment\",
            \"post_id\": \"$POST_ID\",
            \"author_id\": \"$USER_ID\"
        }")
    COMMENT_ID=$(echo "$COMMENT_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    print_result $? "Create comment"

    # Test creating a category
    echo "Creating test category..."
    CATEGORY_RESPONSE=$(curl -s -X POST "$BASE_URL/api/categories" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Test Category",
            "description": "This is a test category"
        }')
    CATEGORY_ID=$(echo "$CATEGORY_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    print_result $? "Create category"

    # Test creating a tag
    echo "Creating test tag..."
    TAG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/tags" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "test-tag"
        }')
    TAG_ID=$(echo "$TAG_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    print_result $? "Create tag"

    # Test linking post to category
    echo "Linking post to category..."
    curl -s -X POST "$BASE_URL/api/posts/$POST_ID/categories" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"category_id\": \"$CATEGORY_ID\"}"
    print_result $? "Link post to category"

    # Test linking post to tag
    echo "Linking post to tag..."
    curl -s -X POST "$BASE_URL/api/posts/$POST_ID/tags" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"tag_id\": \"$TAG_ID\"}"
    print_result $? "Link post to tag"

    # Test creating a media entry
    echo "Creating test media entry..."
    MEDIA_RESPONSE=$(curl -s -X POST "$BASE_URL/api/media" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"url\": \"$URL\",
            \"filename\": \"$FILENAME\",
            \"type\": \"image\",
            \"post_id\": \"$POST_ID\"
        }")
    MEDIA_ID=$(echo "$MEDIA_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    print_result $? "Create media entry"

    # Test creating a setting
    echo "Creating test setting..."
    SETTING_RESPONSE=$(curl -s -X POST "$BASE_URL/api/settings" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "key": "test_setting",
            "value": "test_value",
            "description": "This is a test setting"
        }')
    print_result $? "Create setting"

    # Test creating a menu
    echo "Creating test menu..."
    MENU_RESPONSE=$(curl -s -X POST "$BASE_URL/api/menus" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Test Menu",
            "location": "header"
        }')
    MENU_ID=$(echo "$MENU_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    print_result $? "Create menu"

    # Test creating a menu item
    echo "Creating test menu item..."
    curl -s -X POST "$BASE_URL/api/menus/$MENU_ID/items" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "title": "Test Menu Item",
            "url": "/test",
            "order": 1
        }'
    print_result $? "Create menu item"

    # Test creating a widget
    echo "Creating test widget..."
    curl -s -X POST "$BASE_URL/api/widgets" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "title": "Test Widget",
            "content": "This is a test widget",
            "type": "text",
            "location": "sidebar"
        }'
    print_result $? "Create widget"
}

# Get admin token
echo "Getting admin token..."
TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "admin",
        "password": "pass123"
    }')
TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get admin token${NC}"
    exit 1
fi

echo -e "${GREEN}Got admin token${NC}"

# Copy test files if they exist
if [ -f "../test_files/WIN_20240917_11_28_28_Pro.jpg" ]; then
    cp "../test_files/WIN_20240917_11_28_28_Pro.jpg" "$TEST_FILES_DIR/"
    echo "Copied test image"
else
    echo -e "${RED}Warning: Test image not found${NC}"
fi

# Create a test markdown file
echo "# Test Markdown" > "$TEST_FILES_DIR/test.md"
echo "This is a test markdown file." >> "$TEST_FILES_DIR/test.md"
echo "Created test markdown file"

# Test image upload
test_file_upload "$TEST_FILES_DIR/WIN_20240917_11_28_28_Pro.jpg" "image" "images"

# Test markdown upload
test_file_upload "$TEST_FILES_DIR/test.md" "markdown" "markdown"

# Test database operations
test_database

echo "Tests completed!" 