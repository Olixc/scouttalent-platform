#!/bin/bash

# ScoutTalent API Testing Script
# Tests the complete flow: Register → Login → Create Profile → Update Profile

set -e

BASE_URL_AUTH="http://localhost:8080"
BASE_URL_PROFILE="http://localhost:8081"

echo "🚀 ScoutTalent API Testing"
echo "=========================="
echo ""

# Generate random email for testing
RANDOM_EMAIL="player$(date +%s)@example.com"
PASSWORD="password123"

echo "📝 Step 1: Register a new player"
echo "Email: $RANDOM_EMAIL"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL_AUTH/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$RANDOM_EMAIL\",
    \"password\": \"$PASSWORD\",
    \"role\": \"player\"
  }")

echo "Response: $REGISTER_RESPONSE"
echo ""

# Extract user ID
USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.user.id')
echo "User ID: $USER_ID"
echo ""

echo "🔐 Step 2: Login"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL_AUTH/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$RANDOM_EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "Response: $LOGIN_RESPONSE"
echo ""

# Extract access token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')
echo "Access Token: ${ACCESS_TOKEN:0:50}..."
echo ""

echo "✅ Step 3: Verify authentication"
ME_RESPONSE=$(curl -s -X GET "$BASE_URL_AUTH/api/v1/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "Response: $ME_RESPONSE"
echo ""

echo "👤 Step 4: Create profile"
PROFILE_RESPONSE=$(curl -s -X POST "$BASE_URL_PROFILE/api/v1/profiles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"display_name\": \"John Doe\",
    \"bio\": \"Aspiring professional footballer from Lagos\",
    \"location_country\": \"Nigeria\",
    \"location_city\": \"Lagos\"
  }")

echo "Response: $PROFILE_RESPONSE"
echo ""

# Extract profile ID
PROFILE_ID=$(echo $PROFILE_RESPONSE | jq -r '.id')
echo "Profile ID: $PROFILE_ID"
echo ""

echo "⚽ Step 5: Add player details"
PLAYER_DETAILS_RESPONSE=$(curl -s -X POST "$BASE_URL_PROFILE/api/v1/profiles/$PROFILE_ID/player-details" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"position\": \"forward\",
    \"height_cm\": 180,
    \"weight_kg\": 75,
    \"preferred_foot\": \"right\",
    \"current_team\": \"Lagos FC\"
  }")

echo "Response: $PLAYER_DETAILS_RESPONSE"
echo ""

echo "📊 Step 6: Get complete player profile"
FULL_PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL_PROFILE/api/v1/profiles/$PROFILE_ID/player" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "Response: $FULL_PROFILE_RESPONSE"
echo ""

echo "✨ Step 7: Update profile"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL_PROFILE/api/v1/profiles/$PROFILE_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"bio\": \"Professional footballer with 5 years experience. Specializing in striker position.\"
  }")

echo "Response: $UPDATE_RESPONSE"
echo ""

echo "🎉 All tests completed successfully!"
echo ""
echo "Summary:"
echo "--------"
echo "User ID: $USER_ID"
echo "Profile ID: $PROFILE_ID"
echo "Email: $RANDOM_EMAIL"
echo "Password: $PASSWORD"
echo "Access Token: ${ACCESS_TOKEN:0:50}..."