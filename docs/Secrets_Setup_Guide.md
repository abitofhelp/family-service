# Secrets Setup Guide

## Overview

This document explains the requirements for the `secrets` folder, which contains sensitive configuration files needed for the application to function properly. These files are not included in the repository for security reasons and must be created manually on your development machine.

## Important Notes

1. The `secrets` folder is listed in `.gitignore` and should **never** be committed to the repository
2. All files in this folder are required for the application to function properly
3. You must create these files manually on your development machine

## Required Files

The following files must be created in the `secrets` folder:

### Grafana Security Admin Credentials

1. **gf_security_admin_user.txt**
   - Purpose: Contains the username for Grafana admin access
   - Content: `admin`

2. **gf_security_admin_password.txt**
   - Purpose: Contains the password for Grafana admin access
   - Content: Replace with your secure password
   - Example: `password_here__password_here__password_here__password_here__password_here__password_here__password_here__`

### PostgreSQL Credentials

3. **postgresql_username.txt**
   - Purpose: Contains the username for PostgreSQL database access
   - Content: `postgres`

4. **postgresql_password.txt**
   - Purpose: Contains the password for PostgreSQL database access
   - Content: Replace with your secure password
   - Example: `password_here__password_here__password_here__password_here__password_here__password_here__password_here__`

### MongoDB Credentials

5. **mongo_root_username.txt**
   - Purpose: Contains the username for MongoDB root access
   - Content: `root`

6. **mongo_root_password.txt**
   - Purpose: Contains the password for MongoDB root access
   - Content: Replace with your secure password
   - Example: `password_here__password_here__password_here__password_here__password_here__password_here__password_here__`

### Redis Credentials

7. **redis_password.txt**
   - Purpose: Contains the password for Redis access
   - Content: Replace with your secure password
   - Example: `password_here__password_here__password_here__password_here__password_here__password_here__password_here__`

## Security Recommendations

1. Use strong, unique passwords for each service
2. Do not share these credentials with unauthorized users
3. Consider using a password manager to generate and store these credentials
4. In production environments, consider using a secrets management solution like HashiCorp Vault or AWS Secrets Manager

## How These Secrets Are Used

These secrets are used by the application to connect to various services:

- The Grafana credentials are used to access the Grafana dashboard
- The PostgreSQL credentials are used to connect to the PostgreSQL database
- The MongoDB credentials are used to connect to the MongoDB database
- The Redis credentials are used to connect to the Redis cache

For more information on how these services are configured, please refer to the [Deployment Document](./Deployment_FamilyService.md).