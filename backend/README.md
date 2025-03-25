# Mediation Platform - Backend

- [Mediation Platform - Backend](#mediation-platform---backend)
  - [Dependence](#dependence)
  - [Local Deployment](#local-deployment)
    - [Run Backend API Service](#run-backend-api-service)
    - [API Document](#api-document)
    - [Create Token for Local Testing](#create-token-for-local-testing)
  - [Testing](#testing)

---

## Dependence

- Golang: `1.24.1`
- Swagger fo Golang
  - https://github.com/swaggo/swag
- MongoDB

## Local Deployment

### Run Backend API Service

```bash
make run-app # default app: api-service
```

### API Document

- http://127.0.0.1:8080/swagger/index.html

### Create Token for Local Testing

```bash
make create-token # default TestingUser1
# or
make create-token UID=${FIREBASE_UID}
```

`FIREBASE_UID`: Please refer to [firebase account data](./firebase/emulator_data/auth_export/accounts.json)

## Testing

```bash
make test
```
