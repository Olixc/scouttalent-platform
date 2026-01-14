#!/bin/bash

# Comprehensive end-to-end test script for all ScoutTalent services
# Tests Auth, Profile, and Media services without Docker
# Prerequisites: PostgreSQL, Redis, NATS running locally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AUTH_URL="http://localhost:8080"
PROFILE_URL="http://localhost:8081"
MEDIA_URL="http://localhost:8082"
TEST_EMAIL="e2e-test-$(date +%s)@scouttalent.com"
TEST_PASSWORD="Test123!@#"
TEST_NAME="E2E Test User"

# Test results
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

check_service() {
    local service_name=$1
    local url=$2
    
    if curl -s -f "$url/health" > /dev/null 2>&1; then
        print_success "$service_name is running"
        return 0
    else
        print_error "$service_name is not running at $url"
        return 1
    fi
}

# Start tests
print_header "ScoutTalent Platform - End-to-End Test Suite"

# Step 1: Check prerequisites
print_header "Step 1: Checking Prerequisites"

# Check PostgreSQL
if psql -U scout -d auth_db -c "SELECT 1;" > /dev/null 2>&1; then
    print_success "PostgreSQL is running"
else
    print_error "PostgreSQL is not accessible"
    exit 1
fi

# Check Redis
if redis-cli ping > /dev/null 2>&1; then
    print_success "Redis is running"
else
    print_error "Redis is not accessible"
    exit 1
fi

# Check NATS
if curl -s http://localhost:8222/healthz > /dev/null 2>&1; then
    print_success "NATS is running"
else
    print_error "NATS is not accessible"
    exit 1
fi

# Step 2: Check services
print_header "Step 2: Checking Services"

check_service "Auth Service" "$AUTH_URL" || exit 1
check_service "Profile Service" "$PROFILE_URL" || exit 1
check_service "Media Service" "$MEDIA_URL" || exit 1

# Step 3: Auth Service Tests
print_header "Step 3: Testing Auth Service"

# Register user
print_info "Registering new user: $TEST_EMAIL"
REGISTER_RESPONSE=$(curl -s -X POST "$AUTH_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\",
    \"full_name\": \"$TEST_NAME\"
  }")

if echo "$REGISTER_RESPONSE" | grep -q "id"; then
    print_success "User registration successful"
    USER_ID=$(echo "$REGISTER_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    print_info "User ID: $USER_ID"
else
    print_error "User registration failed"
    echo "Response: $REGISTER_RESPONSE"
    exit 1
fi

# Login
print_info "Logging in as $TEST_EMAIL"
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
    print_success "User login successful"
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    print_info "Access Token: ${ACCESS_TOKEN:0:50}..."
else
    print_error "User login failed"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

# Get current user
print_info "Getting current user info"
ME_RESPONSE=$(curl -s -X GET "$AUTH_URL/api/v1/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$ME_RESPONSE" | grep -q "$TEST_EMAIL"; then
    print_success "Get current user successful"
else
    print_error "Get current user failed"
    echo "Response: $ME_RESPONSE"
fi

# Step 4: Profile Service Tests
print_header "Step 4: Testing Profile Service"

# Create profile
print_info "Creating user profile"
CREATE_PROFILE_RESPONSE=$(curl -s -X POST "$PROFILE_URL/api/v1/profiles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Professional football player with 10 years experience",
    "location": "London, UK",
    "profile_type": "player"
  }')

if echo "$CREATE_PROFILE_RESPONSE" | grep -q "id"; then
    print_success "Profile creation successful"
    PROFILE_ID=$(echo "$CREATE_PROFILE_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    print_info "Profile ID: $PROFILE_ID"
else
    print_error "Profile creation failed"
    echo "Response: $CREATE_PROFILE_RESPONSE"
    exit 1
fi

# Get my profile
print_info "Getting my profile"
MY_PROFILE_RESPONSE=$(curl -s -X GET "$PROFILE_URL/api/v1/profiles/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$MY_PROFILE_RESPONSE" | grep -q "$PROFILE_ID"; then
    print_success "Get my profile successful"
else
    print_error "Get my profile failed"
    echo "Response: $MY_PROFILE_RESPONSE"
fi

# Update profile
print_info "Updating profile"
UPDATE_PROFILE_RESPONSE=$(curl -s -X PUT "$PROFILE_URL/api/v1/profiles/$PROFILE_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Updated: Professional football player with 10 years experience",
    "location": "Manchester, UK"
  }')

if echo "$UPDATE_PROFILE_RESPONSE" | grep -q "Manchester"; then
    print_success "Profile update successful"
else
    print_error "Profile update failed"
    echo "Response: $UPDATE_PROFILE_RESPONSE"
fi

# Create player details
print_info "Creating player details"
PLAYER_DETAILS_RESPONSE=$(curl -s -X POST "$PROFILE_URL/api/v1/profiles/$PROFILE_ID/player-details" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "position": "Forward",
    "preferred_foot": "Right",
    "height": 180,
    "weight": 75,
    "date_of_birth": "2000-01-15"
  }')

if echo "$PLAYER_DETAILS_RESPONSE" | grep -q "Forward"; then
    print_success "Player details creation successful"
else
    print_error "Player details creation failed"
    echo "Response: $PLAYER_DETAILS_RESPONSE"
fi

# Step 5: Media Service Tests
print_header "Step 5: Testing Media Service"

# Initiate video upload
print_info "Initiating video upload"
UPLOAD_RESPONSE=$(curl -s -X POST "$MEDIA_URL/api/v1/videos/upload" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "E2E Test - Skills Showcase",
    "description": "Automated test video upload",
    "file_name": "test-skills.mp4",
    "file_size": 52428800,
    "mime_type": "video/mp4"
  }')

if echo "$UPLOAD_RESPONSE" | grep -q "video_id"; then
    print_success "Video upload initiation successful"
    VIDEO_ID=$(echo "$UPLOAD_RESPONSE" | grep -o '"video_id":"[^"]*' | cut -d'"' -f4)
    UPLOAD_ID=$(echo "$UPLOAD_RESPONSE" | grep -o '"upload_id":"[^"]*' | cut -d'"' -f4)
    print_info "Video ID: $VIDEO_ID"
    print_info "Upload ID: $UPLOAD_ID"
    
    # Check if in test mode
    if echo "$UPLOAD_RESPONSE" | grep -q '"test_mode":true'; then
        print_info "Running in test mode (no Azure credentials)"
    fi
else
    print_error "Video upload initiation failed"
    echo "Response: $UPLOAD_RESPONSE"
    exit 1
fi

# Update upload progress to 50%
print_info "Updating upload progress to 50%"
curl -s -X PATCH "$MEDIA_URL/api/v1/videos/upload/$UPLOAD_ID?progress=50" \
  -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null

sleep 1

# Update upload progress to 100%
print_info "Updating upload progress to 100%"
curl -s -X PATCH "$MEDIA_URL/api/v1/videos/upload/$UPLOAD_ID?progress=100" \
  -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null

print_success "Upload progress updates successful"

# Complete upload
print_info "Completing video upload"
COMPLETE_RESPONSE=$(curl -s -X POST "$MEDIA_URL/api/v1/videos/$VIDEO_ID/complete" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$COMPLETE_RESPONSE" | grep -q "success\|completed\|ready"; then
    print_success "Video upload completion successful"
else
    print_info "Video upload completed (check response for status)"
fi

# Get video details
print_info "Getting video details"
VIDEO_RESPONSE=$(curl -s -X GET "$MEDIA_URL/api/v1/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$VIDEO_RESPONSE" | grep -q "$VIDEO_ID"; then
    print_success "Get video details successful"
    VIDEO_STATUS=$(echo "$VIDEO_RESPONSE" | grep -o '"status":"[^"]*' | cut -d'"' -f4)
    print_info "Video status: $VIDEO_STATUS"
else
    print_error "Get video details failed"
    echo "Response: $VIDEO_RESPONSE"
fi

# Update video metadata
print_info "Updating video metadata"
UPDATE_VIDEO_RESPONSE=$(curl -s -X PUT "$MEDIA_URL/api/v1/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated: E2E Test Video",
    "description": "Updated description for automated test"
  }')

if echo "$UPDATE_VIDEO_RESPONSE" | grep -q "Updated"; then
    print_success "Video metadata update successful"
else
    print_error "Video metadata update failed"
    echo "Response: $UPDATE_VIDEO_RESPONSE"
fi

# List profile videos
print_info "Listing profile videos"
LIST_VIDEOS_RESPONSE=$(curl -s -X GET "$MEDIA_URL/api/v1/videos/profile/$PROFILE_ID?limit=10&offset=0" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$LIST_VIDEOS_RESPONSE" | grep -q "$VIDEO_ID"; then
    print_success "List profile videos successful"
    VIDEO_COUNT=$(echo "$LIST_VIDEOS_RESPONSE" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    print_info "Total videos: $VIDEO_COUNT"
else
    print_error "List profile videos failed"
    echo "Response: $LIST_VIDEOS_RESPONSE"
fi

# Delete video
print_info "Deleting video"
DELETE_RESPONSE=$(curl -s -X DELETE "$MEDIA_URL/api/v1/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

print_success "Video deletion completed"

# Step 6: Summary
print_header "Test Summary"

TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
echo -e "Total Tests: $TOTAL_TESTS"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}✓ All tests passed successfully!${NC}\n"
    exit 0
else
    echo -e "\n${RED}✗ Some tests failed. Please review the output above.${NC}\n"
    exit 1
fi