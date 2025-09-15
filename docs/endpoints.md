# API Endpoints Documentation

This document provides a reference for all the available endpoints in the API-wrkf project.

**Base URL:** `http://localhost:8080`

## 1. Authentication

### Login

-   **Endpoint:** `POST /login`
-   **Description:** Authenticates a user with their email and password to receive a JWT.
-   **Access:** Public
-   **Request Body:**
    ```json
    {
      "correo": "user@example.com",
      "contraseña": "password123"
    }
    ```
-   **Success Response:** `200 OK`
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

---

## 2. Users

### Get User by ID

-   **Endpoint:** `GET /api/users/:id`
-   **Description:** Retrieves the public profile of a specific user.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK`
    ```json
    {
        "ID": 4,
        "Nombre": "Carlos",
        "ApellidoPaterno": "Gomez",
        "ApellidoMaterno": "",
        "Correo": "carlos.gomez@example.com",
        "Role": "user",
        "CreatedAt": "..."
    }
    ```

---

## 3. Projects

### Create Project

-   **Endpoint:** `POST /api/projects`
-   **Description:** Creates a new project. The user making the request is automatically assigned as the creator.
-   **Access:** Authenticated (any valid user)
-   **Request Body:**
    ```json
    {
      "Name": "New Mobile App",
      "Description": "A project to develop a new mobile application."
    }
    ```
-   **Success Response:** `201 Created` with the newly created project object.

### Get All Projects

-   **Endpoint:** `GET /api/projects`
-   **Description:** Retrieves a list of all projects in the system.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK` with an array of project objects.

### Get Project by ID

-   **Endpoint:** `GET /api/projects/:id`
-   **Description:** Retrieves the complete details of a single project, including its creator and a list of its members.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK` with the detailed project object, including `CreatedBy` and `Members` data.

### Update Project

-   **Endpoint:** `PUT /api/projects/:id`
-   **Description:** Updates the details (e.g., name, description) of an existing project.
-   **Access:** Authenticated (Project Creator or Admin only)
-   **Request Body:**
    ```json
    {
      "Name": "New Project Name (Updated)",
      "Description": "An updated description."
    }
    ```
-   **Success Response:** `200 OK` with the updated project object.

### Delete Project

-   **Endpoint:** `DELETE /api/projects/:id`
-   **Description:** Deletes a project and all its associated memberships.
-   **Access:** Authenticated (Project Creator or Admin only)
-   **Success Response:** `204 No Content`

---

## 4. Administration (Admin-Only)

All endpoints in this section require the user to have an `admin` platform role.

### Create Standard User

-   **Endpoint:** `POST /api/admin/users`
-   **Description:** Creates a new standard platform user with the default role of `user`.
-   **Access:** Admin only
-   **Request Body:**
    ```json
    {
      "Nombre": "New",
      "ApellidoPaterno": "User",
      "Correo": "new.user@example.com",
      "Contraseña": "securepassword"
    }
    ```
-   **Success Response:** `201 Created` with the new user object.

### Create Admin User

-   **Endpoint:** `POST /api/admin/users/admin`
-   **Description:** Creates a new administrator user with the platform role of `admin`.
-   **Access:** Admin only
-   **Request Body:** (Same as creating a standard user)
-   **Success Response:** `201 Created` with the new admin user object.

### Add Member to Project

-   **Endpoint:** `POST /api/admin/projects/:id/members`
-   **Description:** Assigns an existing user to a project with a specific project role.
-   **Access:** Admin only
-   **Request Body:**
    ```json
    {
      "userId": 4,
      "role": "team_developer" 
    }
    ```
    *Valid roles are: `scrum_master`, `product_owner`, `team_developer`.*
-   **Success Response:** `201 Created` with the new project membership object.
