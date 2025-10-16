# Admin API

This document outlines the administrative endpoints available in the API.

## Get All Users

This endpoint retrieves a list of all users in the database. This is an admin-only endpoint.

- **URL:** `/api/admin/users`
- **Method:** `GET`
- **Authentication:** Admin JWT Token required.

### Success Response

- **Code:** `200 OK`
- **Content:**

```json
[
    {
        "ID": 1,
        "Nombre": "Admin",
        "ApellidoPaterno": "User",
        "ApellidoMaterno": "",
        "Correo": "admin@example.com",
        "Role": "admin",
        "CreatedAt": "2025-10-16T10:18:20.056-04:00"
    },
    {
        "ID": 2,
        "Nombre": "Test",
        "ApellidoPaterno": "User",
        "ApellidoMaterno": "",
        "Correo": "test@example.com",
        "Role": "user",
        "CreatedAt": "2025-10-16T10:20:00.000-04:00"
    }
]
```

### Error Response

- **Code:** `401 Unauthorized`
  - **Content:** `{"error":"Missing or invalid token"}`
- **Code:** `403 Forbidden`
  - **Content:** `{"error":"Admin access required"}`
- **Code:** `500 Internal Server Error`
  - **Content:** `{"error":"Could not retrieve users"}`
