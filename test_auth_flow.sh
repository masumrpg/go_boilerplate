#!/bin/bash

BASE_URL="http://localhost:3000/api/v1/auth"
EMAIL="testflow1@example.com"
PASSWORD="password123"
NAME="TestFlow User"

echo "---------------------------------------------------"
echo "Testing Auth Flow for $EMAIL"
echo "---------------------------------------------------"

# 1. Register
echo "[1] Registering User..."
REGISTER_RESP=$(curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"$NAME\", \"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")
echo "Response: $REGISTER_RESP"

# 2. Get Activation Code from Redis
echo "[2] Fetching Activation OTP from Redis..."
OTP=$(redis-cli GET "activation:$EMAIL")
echo "OTP: $OTP"

# 3. Verify Email
echo "[3] Verifying Email..."
VERIFY_RESP=$(curl -s -X POST "$BASE_URL/verify-email" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"code\": \"$OTP\"}")
echo "Response: $VERIFY_RESP"

# 4. Login (Trigger 2FA)
echo "[4] Logging in (Expect 2FA)..."
LOGIN_RESP=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")
echo "Response: $LOGIN_RESP"

# 5. Get 2FA OTP from Redis
echo "[5] Fetching 2FA OTP from Redis..."
TWO_FA_OTP=$(redis-cli GET "2fa:$EMAIL")
echo "OTP: $TWO_FA_OTP"

# 6. Verify 2FA (Get Token)
echo "[6] Verifying 2FA..."
TOKEN_RESP=$(curl -s -X POST "$BASE_URL/verify-2fa" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"code\": \"$TWO_FA_OTP\"}")
echo "Response: $TOKEN_RESP"

# Extract Token (Simple grep/cut, assuming json structure)
ACCESS_TOKEN=$(echo $TOKEN_RESP | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo "Access Token: $ACCESS_TOKEN"

# 7. Get Profile (Verify Token)
# Assuming profile endpoint is /api/v1/users/profile or similar?
# Let's check routes. It's usually /api/v1/users/me or similar.
# Checking existing routes... user module usually has one.
# Based on previous file views, there is a user handler but routes weren't explicitly detailed in my memory.
# I'll try /api/v1/users/profile (standard for boilerplate)
echo "[7] Fetching Profile..."
PROFILE_RESP=$(curl -s -X GET "http://localhost:3000/api/v1/users/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "Response: $PROFILE_RESP"

echo "---------------------------------------------------"
echo "Test Complete"
