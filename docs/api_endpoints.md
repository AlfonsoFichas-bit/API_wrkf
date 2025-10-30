# API Endpoints

This document outlines the available API endpoints for the AcademiaSys project management platform.

## Authentication

### POST /login

- **Description:** Authenticates a user and returns a JWT token.
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Response:**
  ```json
  {
    "token": "your.jwt.token"
  }
  ```

### POST /register

- **Description:** Creates a new user account.
- **Request Body:**
  ```json
  {
    "nombre": "John",
    "apellidoPaterno": "Doe",
    "apellidoMaterno": "Smith",
    "correo": "john.doe@example.com",
    "contraseña": "password123"
  }
  ```
- **Response:** The newly created user object.

---

## Evaluations

### POST /api/evaluations

- **Authentication:** Required.
- **Description:** Creates a new evaluation for a task.
- **Request Body:** An `Evaluation` object.
- **Response:** The newly created evaluation object.

### GET /api/evaluations/:id

- **Authentication:** Required.
- **Description:** Retrieves a specific evaluation by its ID.
- **Response:** An `Evaluation` object.

### GET /api/students/:studentId/evaluations

- **Authentication:** Required.
- **Description:** Retrieves all evaluations for a specific student.
- **Response:** An array of `Evaluation` objects.
