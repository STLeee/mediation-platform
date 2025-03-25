# Mediation Platform

- [Mediation Platform](#mediation-platform)
  - [Features](#features)
    - [Frontend](#frontend)
    - [Backend](#backend)
  - [Dependence](#dependence)
  - [Local Development](#local-development)
    - [Setup Firebase Local Emulator](#setup-firebase-local-emulator)
    - [Run All Local Infra Emulators](#run-all-local-infra-emulators)
    - [Develop with Backend](#develop-with-backend)

---

## Features

### Frontend

- [ ] Sign-up & Sign-in
- [ ] Create and List Issues
- [ ] Add and List Comments
- [ ] Notifications

### Backend

- [x] Authentication
- [ ] Issue API
- [ ] Comments API
- [ ] Notifications
- [ ] AI Comment

## Dependence

- Firebase
- [Backend Dependence](./backend/README.md#dependence)

## Local Development

### Setup Firebase Local Emulator

Reference:
- https://firebase.google.com/docs/cli
- https://firebase.google.com/docs/emulator-suite/install_and_configure

1. Install Firebase Cli
2. Login Firebase
    ```bash
    firebase login
    firebase projects:list
    ```
3. Run Emulator
    ```bash
    make run-firebase-emulators
    ```

### Run All Local Infra Emulators

1. Run
    ```bash
    make run-local-infra
    ```

2. Stop
   ```bash
   make stop-local-infra
   ```

### Develop with Backend

Please refer to [Backend document](./backend/README.md#local-deployment)
