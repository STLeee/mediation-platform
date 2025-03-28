# Mediation Platform - Backend

- [Mediation Platform - Backend](#mediation-platform---backend)
  - [Dependence](#dependence)
  - [Local Deployment](#local-deployment)
    - [Run Backend API Service](#run-backend-api-service)
    - [API Document](#api-document)
    - [Generate Token for Local Testing](#generate-token-for-local-testing)
    - [Issues](#issues)
      - [permissions on /data/tls/mongodb-test-keyfile are too open](#permissions-on-datatlsmongodb-test-keyfile-are-too-open)
  - [Testing](#testing)

---

## Dependence

- Golang: `1.24.1`
- Swagger fo Golang
  - https://github.com/swaggo/swag
- MongoDB: `8.0.4`
- Redis: `7.4.2`

## Local Deployment

### Run Backend API Service

```bash
make run-app # default app: api-service
```

### API Document

- http://127.0.0.1:8080/swagger/index.html

### Generate Token for Local Testing

```bash
make generate-token # default TestingUser1
# or
make generate-token UID=${FIREBASE_UID}
```

`FIREBASE_UID`: Please refer to [firebase account data](./firebase/emulator_data/auth_export/accounts.json)

### Issues

#### permissions on /data/tls/mongodb-test-keyfile are too open

```bash
cd mongodb/tls
chmod 600 mongodb-test-keyfile
```

## Testing

```bash
make test
```
