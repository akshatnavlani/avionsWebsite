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

# Function to verify file in Supabase storage
verify_supabase_file() {
    local bucket="$1"
    local filename="$2"
    
    echo "Verifying file in Supabase storage: $bucket/$filename"
    
    # Try to get the file metadata from Supabase
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/storage/$bucket/$filename" \
        -H "Authorization: Bearer $TOKEN")
    
    if [ $? -eq 0 ] && [ -n "$RESPONSE" ]; then
        echo -e "${GREEN}File verified in Supabase storage${NC}"
        return 0
    else
        echo -e "${RED}File not found in Supabase storage${NC}"
        return 1
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
        
        # Verify file in Supabase storage
        verify_supabase_file "$expected_bucket" "$FILENAME"
        return $?
    else
        echo -e "${RED}Upload failed${NC}"
        echo "Response: $RESPONSE"
        return 1
    fi
}

# Function to process markdown content and upload images
process_markdown() {
    local markdown_file="$1"
    local temp_dir="$2"
    
    # Create temp directory for images
    mkdir -p "$temp_dir"
    
    # Extract image references from markdown (handling all common formats)
    # Format 1: ![[filename]]
    # Format 2: ![alt text](image/filename)
    # Format 3: ![alt text](filename)
    # Format 4: <img src="filename" alt="alt text">
    # Format 5: ![alt text](filename "title")
    grep -o '!\[\[.*?\]\]\|!\[.*?\]\(.*?\)\|<img.*?src=".*?".*?>\|!\[.*?\]\(.*?"[^"]*"\)' "$markdown_file" | while read -r line; do
        # Extract image filename (handle all formats)
        if [[ $line =~ !\[\[(.*?)\]\] ]]; then
            # Format 1: ![[filename]]
            img_filename="${BASH_REMATCH[1]}"
        elif [[ $line =~ !\[.*?\]\((image/)?(.*?)(\s+"[^"]*")?\) ]]; then
            # Format 2 & 3 & 5: ![alt text](image/filename) or ![alt text](filename) or ![alt text](filename "title")
            img_filename="${BASH_REMATCH[2]}"
        elif [[ $line =~ <img.*?src="(.*?)".*?> ]]; then
            # Format 4: <img src="filename" alt="alt text">
            img_filename="${BASH_REMATCH[1]}"
        else
            echo -e "${RED}Warning: Could not parse image reference: $line${NC}"
            continue
        fi
        
        echo "Found image reference: $img_filename"
        
        # Try different possible locations for the image file
        if [ -f "./test_files/$img_filename" ]; then
            cp "./test_files/$img_filename" "$temp_dir/"
        elif [ -f "./test_files/image/$img_filename" ]; then
            cp "./test_files/image/$img_filename" "$temp_dir/"
        else
            echo -e "${RED}Warning: Image file not found: $img_filename${NC}"
        fi
    done
}

# Get admin token
echo "Getting admin token..."
TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "password": "pass123"
    }')
TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get admin token${NC}"
    exit 1
fi

echo -e "${GREEN}Got admin token${NC}"

# Create club head member
echo "Creating club head member..."
MEMBER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/members" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Akshat Navlani",
        "position": "Club Head",
        "imageUrl": ""
    }')
MEMBER_ID=$(echo "$MEMBER_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
print_result $? "Create club head member"

# Upload club head profile image
echo "Uploading club head profile image..."
test_file_upload "./test_files/1.jpg" "image" "images"

# Update member with image URL
echo "Updating member with image URL..."
IMAGE_URL=$(curl -s -X POST "$BASE_URL/api/storage/upload" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: multipart/form-data" \
    -F "type=image" \
    -F "file=@./test_files/1.jpg" | grep -o '"url":"[^"]*"' | cut -d'"' -f4)

curl -s -X PUT "$BASE_URL/api/members/$MEMBER_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"Akshat Navlani\",
        \"position\": \"Club Head\",
        \"imageUrl\": \"$IMAGE_URL\"
    }"

# Create a project
echo "Creating project..."
PROJECT_CONTENT=$(cat test_files/test.md | jq -R . | jq -s .)
PROJECT_JSON=$(jq -n \
    --arg title "Drone Racing Project" \
    --arg description "Custom FPV drone racing project" \
    --arg content "$PROJECT_CONTENT" \
    --arg author_id "$MEMBER_ID" \
    --arg status "published" \
    '{
        title: $title,
        description: $description,
        content: ($content | fromjson),
        author_id: $author_id,
        status: $status
    }')

PROJECT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/projects" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$PROJECT_JSON")
PROJECT_ID=$(echo "$PROJECT_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
print_result $? "Create project"

# Create a blog post
echo "Creating blog post..."
BLOG_CONTENT=$(cat test_files/test.md | jq -R . | jq -s .)
BLOG_JSON=$(jq -n \
    --arg title "Drone Racing: A Technical Overview" \
    --arg content "$BLOG_CONTENT" \
    --arg author_id "$MEMBER_ID" \
    --arg status "published" \
    '{
        title: $title,
        content: ($content | fromjson),
        author_id: $author_id,
        status: $status
    }')

BLOG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/posts" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$BLOG_JSON")
BLOG_ID=$(echo "$BLOG_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
print_result $? "Create blog post"

# Process markdown content and upload images
echo "Processing markdown content..."
TEMP_DIR="$TEST_FILES_DIR/temp_images"
process_markdown "test_files/test.md" "$TEMP_DIR"

# Upload images from markdown
for img in "$TEMP_DIR"/*; do
    if [ -f "$img" ]; then
        echo "Uploading image from markdown: $(basename \"$img\")"
        test_file_upload "$img" "image" "images"
    fi
done

# Verify all uploaded files in Supabase
echo "Verifying all uploaded files in Supabase storage..."
verify_supabase_file "images" "1.jpg"
verify_supabase_file "images" "2.jpg"

# Clean up
rm -rf "$TEMP_DIR"

echo "Tests completed!" 