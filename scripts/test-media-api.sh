#!/bin/bash

# Test script for Media Service API
# Prerequisites: 
# 1. All services running (docker-compose up)
# 2. Auth service accessible at localhost:8080
# 3. Media service accessible at localhost:8082

set -e

BASE_URL="http://localhost:8082"
AUTH_URL="http://localhost:8080"

echo "=== ScoutTalent Media Service API Test ==="
echo ""

# Step 1: Register a user
echo "1. Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST "$AUTH_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@test.com",
    "password": "Test123!@#",
    "full_name": "Test Player"
  }')

echo "Register Response: $REGISTER_RESPONSE"
echo ""

# Step 2: Login to get JWT token
echo "2. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@test.com",
    "password": "Test123!@#"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Failed to get access token"
  echo "Login Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "Access Token: ${TOKEN:0:50}..."
echo ""

# Step 3: Check Media Service health
echo "3. Checking Media Service health..."
HEALTH_RESPONSE=$(curl -s "$BASE_URL/health")
echo "Health Response: $HEALTH_RESPONSE"
echo ""

# Step 4: Initiate video upload
echo "4. Initiating video upload..."
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/videos/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Football Skills Showcase",
    "description": "Demonstrating dribbling and shooting techniques",
    "file_name": "skills_demo.mp4",
    "file_size": 52428800,
    "mime_type": "video/mp4"
  }')

echo "Upload Response: $UPLOAD_RESPONSE"

VIDEO_ID=$(echo $UPLOAD_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
UPLOAD_ID=$(echo $UPLOAD_RESPONSE | grep -o '"id":"[^"]*' | tail -1 | cut -d'"' -f4)

if [ -z "$VIDEO_ID" ]; then
  echo "Failed to initiate upload"
  exit 1
fi

echo "Video ID: $VIDEO_ID"
echo "Upload ID: $UPLOAD_ID"
echo ""

# Step 5: Update upload progress
echo "5. Updating upload progress to 50%..."
curl -s -X PATCH "$BASE_URL/api/v1/videos/upload/$UPLOAD_ID?progress=50" \
  -H "Authorization: Bearer $TOKEN"
echo ""

sleep 1

echo "6. Updating upload progress to 100%..."
curl -s -X PATCH "$BASE_URL/api/v1/videos/upload/$UPLOAD_ID?progress=100" \
  -H "Authorization: Bearer $TOKEN"
echo ""

# Step 7: Complete upload
echo "7. Completing upload..."
COMPLETE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/videos/$VIDEO_ID/complete" \
  -H "Authorization: Bearer $TOKEN")
echo "Complete Response: $COMPLETE_RESPONSE"
echo ""

# Step 8: Get video details
echo "8. Getting video details..."
VIDEO_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "Video Details: $VIDEO_RESPONSE"
echo ""

# Step 9: Update video metadata
echo "9. Updating video metadata..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/v1/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated: My Football Skills",
    "description": "Updated description with more details"
  }')
echo "Update Response: $UPDATE_RESPONSE"
echo ""

# Step 10: List profile videos
echo "10. Listing profile videos..."
# Extract profile_id from token claims (simplified - in production use proper JWT parsing)
PROFILE_ID=$(echo $TOKEN | cut -d'.' -f2 | base64 -d 2>/dev/null | grep -o '"profile_id":"[^"]*' | cut -d'"' -f4 || echo "")

if [ ! -z "$PROFILE_ID" ]; then
  LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/videos/profile/$PROFILE_ID?limit=10&offset=0" \
    -H "Authorization: Bearer $TOKEN")
  echo "Profile Videos: $LIST_RESPONSE"
else
  echo "Could not extract profile_id from token"
fi
echo ""

# Step 11: Delete video
echo "11. Deleting video..."
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "Delete Response: $DELETE_RESPONSE"
echo ""

echo "=== Media Service API Test Complete ==="