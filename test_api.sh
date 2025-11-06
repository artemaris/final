#!/bin/bash

# GophKeeper API Test Script
# This script tests the API endpoints

BASE_URL="http://localhost:8080"

echo "========================================"
echo "GophKeeper API Tests"
echo "========================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test 1: Health check (if exists)
echo "1. Testing server health..."
curl -s "$BASE_URL/health" || echo "Health endpoint not implemented"
echo ""

# Test 2: Register a new user
echo "2. Registering a new user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "$REGISTER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$REGISTER_RESPONSE"
echo ""

# Extract token from response
TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get token. Attempting login...${NC}"
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
      -H "Content-Type: application/json" \
      -d '{
        "username": "testuser",
        "password": "password123"
      }')
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Authentication failed. Cannot continue tests.${NC}"
    exit 1
fi

echo -e "${GREEN}Authentication successful!${NC}"
echo ""

# Test 3: List entries (should be empty)
echo "3. Listing all entries..."
curl -s -X GET "$BASE_URL/api/entries" \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo ""

# Test 4: Create a login entry
echo "4. Creating a login entry..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/entries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "login",
    "title": "GitHub",
    "data": "encrypted_data_here",
    "metadata": "{\"site\":\"github.com\"}"
  }')

echo "$CREATE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CREATE_RESPONSE"
echo ""

# Extract entry ID
ENTRY_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

# Test 5: Get specific entry
if [ -n "$ENTRY_ID" ]; then
    echo "5. Getting entry by ID: $ENTRY_ID..."
    curl -s -X GET "$BASE_URL/api/entries/$ENTRY_ID" \
      -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
    echo ""
fi

# Test 6: Update entry
if [ -n "$ENTRY_ID" ]; then
    echo "6. Updating entry..."
    curl -s -X PUT "$BASE_URL/api/entries/$ENTRY_ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "GitHub Updated",
        "data": "updated_encrypted_data",
        "metadata": "{\"site\":\"github.com\",\"updated\":true}",
        "version": 1
      }' | python3 -m json.tool
    echo ""
fi

# Test 7: Sync entries
echo "7. Syncing entries..."
curl -s -X GET "$BASE_URL/api/sync?since=0" \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo ""

echo -e "${GREEN}Tests completed!${NC}"
