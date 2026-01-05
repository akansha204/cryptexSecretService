# CryptexSecretService

Cryptex Secret Service is a backend microservice that provides secure, encrypted, versioned secret storage similar to HashiCorp Vault or AWS Secrets Manager.
It manages Projects, Secrets, Access Policies, and Audit Logs.
This service is built using Go (Fiber + Bun ORM) and PostgreSQL.

Although Cryptex follows a microservices architecture, all microservices share the same database during development. This simplifies schema evolution and ensures consistency while keeping the services independently structured.

---
## Features
### Project Management
Users can create and manage isolated project namespaces.
Each project can contain multiple secrets.

Supported operations:

-> Create project<br>
-> Get project by ID<br>
-> List all projects for a user<br>
-> Update project details<br>
-> Soft delete a project<br>
-> Automatic purge of deleted projects after X days<br>

### Secret Management
Each secret belongs to a project and supports automatic encryption, versioning, TTL-based expiration, and revocation.

Supported operations:

-> Create a secret (auto-encryption and versioning)<br>
-> Retrieve a secret (auto-decryption)<br>
-> Update secret (new version if value changes)<br>
-> Revoke secret<br>
-> Soft delete secret<br>
-> Auto-purge deleted secrets after the retention window<br>

### Audit Logging
Every action affecting Projects or Secrets is recorded in an immutable audit log entry.

The system writes audit logs for:

-> Project creation, updation and deletion.<br>
-> Project creation, updation, deletion and revocation.<br>

---
## API Endpoint Examples
### **POST** `/api/projects`
```bash
Authorization: Bearer <JWT>
```
Create Project

**Request:**
```json
{
  "name": "Go-cli project",
  "description": "music can be played from terminal"
}
```
Create Secret
### **POST** `/api/projects/:projectId/secrets`
The `ttl` (Time-To-Live) value in the JSON body must be specified in **days**. If you won't provide the ttl in req body it expires never.

```json
{
  "name": "SPOTIFY_KEYS",
  "value": "your-spotify-key",
  "ttl": 30 
}
```
Retrieve Secret
### **GET** `/api/projects/:projectId/secrets/:secretId`
Returns decrypted secret value only if:

-> Project belongs to the user<br>
-> Secret is not deleted<br>
-> Secret is not expired<br>
-> Secret is not revoked<br>

Delete Secret (Soft Delete)
### **DELETE** `/api/projects/:projectId/secrets/:secretId`
Logs the deletion event and marks the secret as deleted.

---
## Database Schema Overview
Below is the detailed ER diagram representing the database structure for Cryptex Secret Service. This Service only has **Project**,**Secrets** and **audit** entities.
<img width="1411" height="816" alt="image" src="https://github.com/user-attachments/assets/110d1320-38e1-4ac3-b00f-701705a7470f" />

---
## Single Database for All Microservices
Cryptex uses one shared PostgreSQL database for all microservices.
This approach was chosen intentionally due to early-stage development.




