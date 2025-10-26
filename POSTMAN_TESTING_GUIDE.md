# 🚀 Complete Postman Testing Guide for Go Microservices

## 📋 Table of Contents
1. [Setup & Prerequisites](#setup--prerequisites)
2. [Service Overview & URLs](#service-overview--urls)
3. [Health Check Tests](#health-check-tests)
4. [Broker Service Tests](#broker-service-tests)
5. [Authentication Service Tests](#authentication-service-tests)
6. [Logger Service Tests](#logger-service-tests)
7. [Mail Service Tests](#mail-service-tests)
8. [Advanced Testing Scenarios](#advanced-testing-scenarios)
9. [Troubleshooting](#troubleshooting)
10. [Postman Collection Export](#postman-collection-export)

---

## 🔧 Setup & Prerequisites

### 1. Install Postman
Download from: https://www.postman.com/downloads/

### 2. Start All Services
```powershell
cd "D:\internship assignment\Go\udemy-micro\project"
make up_build
```

### 3. Verify Services are Running
```powershell
docker ps
```

You should see these containers:
- `project-broker-service-1` (Port 8080)
- `project-authentication-service-1` (Port 8081)
- `project-logger-service-1` (Internal only)
- `project-mailer-service-1` (Internal only)
- `project-postgres-1` (Port 5432)
- `project-mongo-1` (Port 27017)
- `project-mailhog-1` (Ports 1025, 8025)

---

## 🌐 Service Overview & URLs

| Service | Port | External Access | Internal Docker Name |
|---------|------|-----------------|---------------------|
| **Broker** (API Gateway) | 8080 | ✅ `http://localhost:8080` | `http://broker-service` |
| **Authentication** | 8081 | ✅ `http://localhost:8081` | `http://authentication-service` |
| **Logger** | - | ❌ Internal only | `http://logger-service` |
| **Mailer** | - | ❌ Internal only | `http://mailer-service` |
| **PostgreSQL** | 5432 | ✅ `localhost:5432` | `postgres` |
| **MongoDB** | 27017 | ✅ `localhost:27017` | `mongo` |
| **MailHog UI** | 8025 | ✅ `http://localhost:8025` | - |

**Important Notes:**
- Logger and Mailer services are **NOT exposed** to the host - they're only accessible from within the Docker network
- All requests to Logger/Mailer **must go through the Broker** service
- MailHog is a fake SMTP server for testing emails

---

## 🏥 Health Check Tests

Test if services are responding (all services have heartbeat middleware).

### 1. Broker Health Check

**Request:**
```
Method: GET
URL: http://localhost:8080/ping
```

**Expected Response:**
```
Status: 200 OK
Body: . (single dot)
```

---

### 2. Authentication Health Check

**Request:**
```
Method: GET
URL: http://localhost:8081/ping
```

**Expected Response:**
```
Status: 200 OK
Body: . (single dot)
```

---

## 🔑 Broker Service Tests

The Broker is your **API Gateway** - all external requests go through it.

### Test 1: Broker Root Endpoint

**Request:**
```
Method: POST
URL: http://localhost:8080/
Headers: (none required)
Body: (none)
```

**Expected Response:**
```json
{
    "error": false,
    "message": "Hit the broker"
}
```

**Status Code:** `200 OK`

---

### Test 2: Authentication via Broker

**Request:**
```
Method: POST
URL: http://localhost:8080/handle
Headers:
  Content-Type: application/json
Body (raw JSON):
{
    "action": "auth",
    "auth": {
        "email": "admin@example.com",
        "password": "verysecret"
    }
}
```

**Expected Response (Success):**
```json
{
    "error": false,
    "message": "Authenticated!",
    "data": {
        "id": 1,
        "email": "admin@example.com",
        "first_name": "Admin",
        "last_name": "User",
        "user_active": 1,
        "created_at": "2025-10-18T00:00:00Z",
        "updated_at": "2025-10-18T00:00:00Z"
    }
}
```

**Status Code:** `202 Accepted`

**Expected Response (Invalid Credentials):**
```json
{
    "error": true,
    "message": "invalid credentials"
}
```

**Status Code:** `401 Unauthorized`

---

### Test 3: Logging via Broker

**Request:**
```
Method: POST
URL: http://localhost:8080/handle
Headers:
  Content-Type: application/json
Body (raw JSON):
{
    "action": "log",
    "log": {
        "name": "test-event",
        "data": "This is a test log entry from Postman"
    }
}
```

**Expected Response:**
```json
{
    "error": false,
    "message": "logged"
}
```

**Status Code:** `202 Accepted`

**Verify in MongoDB Compass:**
1. Open MongoDB Compass
2. Connect to: `mongodb://admin:password@localhost:27017/logs?authSource=admin`
3. Navigate to: `logs` database → `logs` collection
4. You should see your log entry with fields: `name`, `data`, `created_at`

---

### Test 4: Send Mail via Broker

**Request:**
```
Method: POST
URL: http://localhost:8080/handle
Headers:
  Content-Type: application/json
Body (raw JSON):
{
    "action": "mail",
    "mail": {
        "from": "test@example.com",
        "to": "recipient@example.com",
        "subject": "Test Email from Postman",
        "message": "This is a test email sent via the broker service."
    }
}
```

**Expected Response:**
```json
{
    "error": false,
    "message": "Message sent to recipient@example.com"
}
```

**Status Code:** `202 Accepted`

**Verify Email Sent:**
1. Open browser: `http://localhost:8025` (MailHog UI)
2. You should see your test email in the inbox
3. Click to view the full email content

---

### Test 5: Invalid Action (Error Test)

**Request:**
```
Method: POST
URL: http://localhost:8080/handle
Headers:
  Content-Type: application/json
Body (raw JSON):
{
    "action": "invalid-action"
}
```

**Expected Response:**
```json
{
    "error": true,
    "message": "unknown action"
}
```

**Status Code:** `400 Bad Request`

---

## 🔐 Authentication Service Tests

Direct tests to the Authentication service (bypassing broker).

### Test 1: Direct Authentication (Valid)

**Request:**
```
Method: POST
URL: http://localhost:8081/authenticate
Headers:
  Content-Type: application/json
Body (raw JSON):
{
    "email": "admin@example.com",
    "password": "verysecret"
}
```

**Expected Response:**
```json
{
    "error": false,
    "message": "Logged in user admin@example.com",
    "data": {
        "id": 1,
        "email": "admin@example.com",
        "first_name": "Admin",
        "last_name": "User",
        "user_active": 1,
        "created_at": "2025-10-18T00:00:00Z",
        "updated_at": "2025-10-18T00:00:00Z"
    }
}
```

**Status Code:** `202 Accepted`

---

### Test 2: Direct Authentication (Invalid)

**Request:**
```
Method: POST
URL: http://localhost:8081/authenticate
Headers:
  Content-Type: application/json
Body (raw JSON):
{
    "email": "wrong@example.com",
    "password": "wrongpassword"
}
```

**Expected Response:**
```json
{
    "error": true,
    "message": "invalid credentials"
}
```

**Status Code:** `400 Bad Request`

---

## 📝 Logger Service Tests

**Note:** Logger service is NOT exposed externally. You **must** use the broker to test it.

### Test: Log Entry via Broker

See [Test 3: Logging via Broker](#test-3-logging-via-broker) above.

**Additional Test Cases:**

#### Different Log Types
```json
{
    "action": "log",
    "log": {
        "name": "error-log",
        "data": "Critical error occurred in payment processing"
    }
}
```

```json
{
    "action": "log",
    "log": {
        "name": "info-log",
        "data": "User logged in successfully"
    }
}
```

```json
{
    "action": "log",
    "log": {
        "name": "audit-log",
        "data": "{\"user_id\": 123, \"action\": \"delete_record\", \"record_id\": 456}"
    }
}
```

---

## 📧 Mail Service Tests

**Note:** Mailer service is NOT exposed externally. You **must** use the broker to test it.

### Test: Send Various Email Types

#### Test 1: Welcome Email
```json
{
    "action": "mail",
    "mail": {
        "from": "noreply@myapp.com",
        "to": "newuser@example.com",
        "subject": "Welcome to Our Service!",
        "message": "Thank you for signing up. We're excited to have you on board!"
    }
}
```

#### Test 2: Password Reset Email
```json
{
    "action": "mail",
    "mail": {
        "from": "security@myapp.com",
        "to": "user@example.com",
        "subject": "Password Reset Request",
        "message": "Click here to reset your password: http://example.com/reset?token=abc123"
    }
}
```

#### Test 3: Notification Email
```json
{
    "action": "mail",
    "mail": {
        "from": "notifications@myapp.com",
        "to": "user@example.com",
        "subject": "New Message Received",
        "message": "You have a new message from John Doe. Login to view it."
    }
}
```

**View All Emails:**
- Open: `http://localhost:8025`
- MailHog shows all sent emails in a web interface
- No real emails are sent (it's a fake SMTP server for testing)

---

## 🧪 Advanced Testing Scenarios

### Scenario 1: Complete User Workflow

**Step 1: Authenticate**
```json
POST http://localhost:8080/handle
{
    "action": "auth",
    "auth": {
        "email": "admin@example.com",
        "password": "verysecret"
    }
}
```

**Step 2: Log the Authentication**
```json
POST http://localhost:8080/handle
{
    "action": "log",
    "log": {
        "name": "user-login",
        "data": "User admin@example.com logged in successfully"
    }
}
```

**Step 3: Send Welcome Email**
```json
POST http://localhost:8080/handle
{
    "action": "mail",
    "mail": {
        "from": "system@myapp.com",
        "to": "admin@example.com",
        "subject": "Login Notification",
        "message": "You have successfully logged in to your account."
    }
}
```

---

### Scenario 2: Error Handling Test

**Test Missing Required Fields:**

```json
POST http://localhost:8080/handle
{
    "action": "auth",
    "auth": {
        "email": "",
        "password": ""
    }
}
```

**Test Malformed JSON:**

```json
POST http://localhost:8080/handle
{
    "action": "log",
    "log": {
        "name": "test"
        // Missing comma, invalid JSON
        "data": "test"
    }
}
```

---

### Scenario 3: Load Testing

Use Postman's **Collection Runner** to send multiple requests:

1. Create a collection with all your tests
2. Click "Run Collection"
3. Set iterations: 100
4. Set delay: 100ms
5. Run to test performance

---

## 🔍 Troubleshooting

### Problem: Connection Refused

**Symptoms:**
```
Error: connect ECONNREFUSED 127.0.0.1:8080
```

**Solutions:**
1. Check if containers are running:
   ```powershell
   docker ps
   ```

2. Check logs:
   ```powershell
   docker-compose logs broker-service
   docker-compose logs authentication-service
   ```

3. Restart services:
   ```powershell
   make down
   make up_build
   ```

---

### Problem: "error calling auth service"

**Symptoms:**
```json
{
    "error": true,
    "message": "error calling auth service"
}
```

**Solutions:**
1. Check if auth service is running:
   ```powershell
   docker ps | findstr authentication
   ```

2. Test auth service directly:
   ```
   GET http://localhost:8081/ping
   ```

3. Check if database is accessible:
   ```powershell
   docker-compose logs postgres
   ```

4. Verify DSN connection string in `docker-compose.yml`

---

### Problem: Logs Not Appearing in MongoDB

**Symptoms:**
- Request succeeds but no logs in MongoDB

**Solutions:**
1. Check logger service logs:
   ```powershell
   docker-compose logs logger-service
   ```

2. Verify MongoDB connection:
   ```powershell
   docker-compose logs mongo
   ```

3. Check MongoDB Compass connection:
   ```
   mongodb://admin:password@localhost:27017/logs?authSource=admin
   ```

4. Ensure `logs` database and `logs` collection exist

---

### Problem: Emails Not Showing in MailHog

**Symptoms:**
- Mail request succeeds but no email in MailHog

**Solutions:**
1. Check if MailHog is running:
   ```powershell
   docker ps | findstr mailhog
   ```

2. Access MailHog UI:
   ```
   http://localhost:8025
   ```

3. Check mailer service logs:
   ```powershell
   docker-compose logs mailer-service
   ```

4. Verify environment variables in `docker-compose.yml`:
   ```yaml
   MAIL_HOST: mailhog
   MAIL_PORT: 1025
   ```

---

## 📦 Postman Collection Export

### Create a Collection

1. **Open Postman**
2. Click **"New Collection"**
3. Name it: `Go Microservices - Complete Tests`

### Add Folders

Create these folders in your collection:
- `Health Checks`
- `Broker Service`
- `Authentication Service`
- `Logger Service (via Broker)`
- `Mail Service (via Broker)`
- `Advanced Workflows`

### Add Requests to Each Folder

**Example Structure:**
```
📁 Go Microservices - Complete Tests
  📁 Health Checks
    ➤ Broker Health Check (GET /ping)
    ➤ Auth Health Check (GET /ping)
  📁 Broker Service
    ➤ Broker Root (POST /)
    ➤ Auth via Broker (POST /handle)
    ➤ Log via Broker (POST /handle)
    ➤ Mail via Broker (POST /handle)
  📁 Authentication Service
    ➤ Direct Auth - Valid (POST /authenticate)
    ➤ Direct Auth - Invalid (POST /authenticate)
  📁 Logger Service (via Broker)
    ➤ Info Log (POST /handle)
    ➤ Error Log (POST /handle)
    ➤ Audit Log (POST /handle)
  📁 Mail Service (via Broker)
    ➤ Welcome Email (POST /handle)
    ➤ Password Reset (POST /handle)
  📁 Advanced Workflows
    ➤ Complete User Flow (multi-request)
```

### Export Collection

1. Click **three dots (...)** next to collection name
2. Select **"Export"**
3. Choose **"Collection v2.1"**
4. Save as: `Go-Microservices-Tests.postman_collection.json`

### Import Collection (Later Use)

1. Click **"Import"** button
2. Drag and drop the JSON file
3. All requests restored with proper configuration

---

## 🎯 Quick Reference Card

### All Endpoints Summary

| Service | Method | Endpoint | Body Required | Purpose |
|---------|--------|----------|---------------|---------|
| Broker | GET | `/ping` | No | Health check |
| Broker | POST | `/` | No | Root endpoint test |
| Broker | POST | `/handle` | Yes | Route to services |
| Auth | GET | `/ping` | No | Health check |
| Auth | POST | `/authenticate` | Yes | Login user |

### Common Request Bodies

**Authentication:**
```json
{
    "action": "auth",
    "auth": {
        "email": "admin@example.com",
        "password": "verysecret"
    }
}
```

**Logging:**
```json
{
    "action": "log",
    "log": {
        "name": "event-name",
        "data": "log message"
    }
}
```

**Mail:**
```json
{
    "action": "mail",
    "mail": {
        "from": "sender@example.com",
        "to": "recipient@example.com",
        "subject": "Subject Line",
        "message": "Email body"
    }
}
```

---

## 📝 Testing Checklist

Before considering your microservices "tested", verify:

- [ ] All health checks return `200 OK`
- [ ] Broker root endpoint works
- [ ] Authentication works with valid credentials
- [ ] Authentication rejects invalid credentials
- [ ] Logs appear in MongoDB
- [ ] Emails appear in MailHog
- [ ] Invalid actions return proper errors
- [ ] All services restart correctly after `make down && make up`
- [ ] Front-end can communicate with broker
- [ ] Docker logs show no critical errors

---

## 🎓 Learning Outcomes

By completing this guide, you've learned:

✅ How to test microservices architecture  
✅ API Gateway pattern (Broker service)  
✅ Service-to-service communication  
✅ Health check endpoints  
✅ Error handling in distributed systems  
✅ Database integration (PostgreSQL, MongoDB)  
✅ Email testing with MailHog  
✅ Docker networking between services  
✅ Postman advanced features  
✅ RESTful API best practices  

---

## 📚 Additional Resources

- [Postman Documentation](https://learning.postman.com/)
- [Go Chi Router](https://github.com/go-chi/chi)
- [Docker Compose Networking](https://docs.docker.com/compose/networking/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [PostgreSQL with Go](https://www.postgresql.org/docs/)
- [MailHog Documentation](https://github.com/mailhog/MailHog)

---

## 🤝 Contributing

If you find issues or want to add more test cases:
1. Document the test clearly
2. Include expected vs actual results
3. Add troubleshooting steps if needed

---

**Happy Testing! 🚀**

---

*Last Updated: October 26, 2025*  
*Repository: go_microservices-*  
*Author: Your Go Microservices Tutorial*
