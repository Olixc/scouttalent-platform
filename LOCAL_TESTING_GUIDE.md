# Local Testing Guide - ScoutTalent Platform

This guide will walk you through testing the Auth Service and Profile Service on your local machine.

## Prerequisites

Before starting, ensure you have the following installed:

- **Docker Desktop** (version 20.10+)
- **Go** (version 1.23+)
- **Make** (usually pre-installed on macOS/Linux, Windows users can use WSL or install via Chocolatey)
- **curl** (for API testing)
- **jq** (optional, for pretty JSON output)

Verify installations:
```bash
docker --version
docker-compose --version
go version
make --version
```

## Step 1: Clone the Repository

```bash
git clone git@github.com:Olixc/scouttalent-platform.git
cd scouttalent-platform
```

## Step 2: Start Infrastructure Services

Start PostgreSQL, Redis, and NATS using Docker Compose:

```bash
docker-compose up -d
```

Verify all services are running:
```bash
docker-compose ps
```

You should see:
- `postgres` - Running on port 5432
- `redis` - Running on port 6379
- `nats` - Running on port 4222

**Wait 10-15 seconds** for PostgreSQL to fully initialize before proceeding.

## Step 3: Run Database Migrations

Apply migrations for both services:

```bash
# Auth Service migrations
make migrate-up-auth

# Profile Service migrations
make migrate-up-profile
```

Verify migrations succeeded:
```bash
# Check Auth Service database
docker exec -it scouttalent-postgres psql -U postgres -d auth_db -c "\dt"

# Check Profile Service database
docker exec -it scouttalent-postgres psql -U postgres -d profile_db -c "\dt"
```

You should see the `users` table in auth_db and `profiles` table in profile_db.

## Step 4: Start the Services

Open **two separate terminal windows** in the project root directory.

**Terminal 1 - Auth Service:**
```bash
make dev-auth
```

You should see:
```
2025/01/13 16:00:00 INFO Starting Auth Service server=:8080
2025/01/13 16:00:00 INFO Connected to database service=auth
```

**Terminal 2 - Profile Service:**
```bash
make dev-profile
```

You should see:
```
2025/01/13 16:00:05 INFO Starting Profile Service server=:8081
2025/01/13 16:00:05 INFO Connected to database service=profile
```

## Step 5: Run Automated Tests

Open a **third terminal window** and run the automated test script:

```bash
bash scripts/test-api.sh
```

The script will:
1. ✅ Register a new player
2. ✅ Login and receive JWT token
3. ✅ Create a player profile
4. ✅ Retrieve the profile
5. ✅ Update the profile
6. ✅ Test authentication (access protected endpoint)

**Expected Output:**
```
🚀 ScoutTalent API Testing
==========================

📝 Step 1: Register a new player
Email: player1234567890@example.com
Response: {"id":"...","email":"player...","role":"player","created_at":"..."}
✅ Registration successful

🔐 Step 2: Login
Response: {"token":"eyJhbGc...","user":{"id":"...","email":"...","role":"player"}}
✅ Login successful

👤 Step 3: Create player profile
Response: {"id":"...","user_id":"...","full_name":"John Doe","position":"Forward",...}
✅ Profile created

📖 Step 4: Get profile
Response: {"id":"...","user_id":"...","full_name":"John Doe",...}
✅ Profile retrieved

✏️ Step 5: Update profile
Response: {"id":"...","user_id":"...","full_name":"John Doe","position":"Midfielder",...}
✅ Profile updated

🔒 Step 6: Test authentication
Response: {"id":"...","user_id":"...","full_name":"John Doe",...}
✅ Authentication working

🎉 All tests passed!
```

## Step 6: Manual API Testing

You can also test the APIs manually using curl:

### Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "role": "player"
  }' | jq
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }' | jq
```

**Save the token from the response** for subsequent requests.

### Create Profile (Replace YOUR_TOKEN_HERE)

```bash
curl -X POST http://localhost:8081/api/v1/profiles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "full_name": "Test Player",
    "date_of_birth": "2000-01-15",
    "position": "Forward",
    "height": 180,
    "weight": 75,
    "preferred_foot": "right"
  }' | jq
```

### Get Profile (Replace YOUR_TOKEN_HERE)

```bash
curl -X GET http://localhost:8081/api/v1/profiles/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
```

### Update Profile (Replace YOUR_TOKEN_HERE)

```bash
curl -X PUT http://localhost:8081/api/v1/profiles/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "position": "Midfielder",
    "bio": "Experienced midfielder with strong passing skills"
  }' | jq
```

## Step 7: Verify Database Data

Check the data directly in PostgreSQL:

```bash
# View registered users
docker exec -it scouttalent-postgres psql -U postgres -d auth_db -c "SELECT id, email, role, created_at FROM users;"

# View profiles
docker exec -it scouttalent-postgres psql -U postgres -d profile_db -c "SELECT id, user_id, full_name, position, completion_score FROM profiles;"
```

## Troubleshooting

### Issue: "Connection refused" errors

**Solution:** Ensure Docker services are running:
```bash
docker-compose ps
docker-compose logs postgres
```

### Issue: "Database does not exist"

**Solution:** Run migrations:
```bash
make migrate-up-auth
make migrate-up-profile
```

### Issue: "Port already in use"

**Solution:** Check if another service is using ports 8080 or 8081:
```bash
# macOS/Linux
lsof -i :8080
lsof -i :8081

# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :8081
```

Kill the conflicting process or change the port in the service configuration.

### Issue: "Invalid token" errors

**Solution:** Ensure you're using a fresh token from the login response. Tokens expire after 24 hours.

### Issue: Migration fails with "relation already exists"

**Solution:** Reset the database:
```bash
make migrate-down-auth
make migrate-down-profile
make migrate-up-auth
make migrate-up-profile
```

## Stopping the Services

### Stop the Go services
Press `Ctrl+C` in each terminal running the services.

### Stop Docker containers
```bash
docker-compose down
```

### Stop and remove all data (clean slate)
```bash
docker-compose down -v
```

## Next Steps

Once all tests pass successfully:

1. ✅ **Auth Service** - Fully tested and working
2. ✅ **Profile Service** - Fully tested and working
3. 🚧 **Media Service** - Ready to build next (video uploads, AI moderation)

## Test Results Checklist

Before moving to the next service, verify:

- [ ] Docker containers are running (postgres, redis, nats)
- [ ] Database migrations applied successfully
- [ ] Auth Service starts without errors on port 8080
- [ ] Profile Service starts without errors on port 8081
- [ ] User registration works
- [ ] User login returns valid JWT token
- [ ] Profile creation works with authentication
- [ ] Profile retrieval works
- [ ] Profile update works
- [ ] Protected endpoints reject requests without valid tokens
- [ ] Automated test script passes all checks

## Support

If you encounter any issues not covered in this guide:

1. Check service logs in the terminal windows
2. Check Docker logs: `docker-compose logs -f`
3. Verify environment variables in `.env` files (if any)
4. Review the code in `services/auth-service/` and `services/profile-service/`

---

**Ready to test?** Follow the steps above and report back with your results. Once everything passes, we'll move on to building the Media Service! 🚀