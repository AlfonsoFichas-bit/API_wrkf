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

---

## 3. Projects

### Create Project

-   **Endpoint:** `POST /api/projects`
-   **Description:** Creates a new project. The user making the request is automatically assigned as the creator.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `201 Created`

### Get All Projects

-   **Endpoint:** `GET /api/projects`
-   **Description:** Retrieves a list of all projects in the system.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK`

### Get Project by ID

-   **Endpoint:** `GET /api/projects/:id`
-   **Description:** Retrieves the complete details of a single project, including its creator and members.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK`

### Update Project

-   **Endpoint:** `PUT /api/projects/:id`
-   **Description:** Updates the details of an existing project.
-   **Access:** Authenticated (Project Creator or Admin only)
-   **Success Response:** `200 OK`

### Delete Project

-   **Endpoint:** `DELETE /api/projects/:id`
-   **Description:** Deletes a project and all its associated dependencies (members, sprints, user stories, tasks).
-   **Access:** Authenticated (Project Creator or Admin only)
-   **Success Response:** `204 No Content`

---

## 4. User Stories (Product Backlog)

### Create User Story

-   **Endpoint:** `POST /api/projects/:id/userstories`
-   **Description:** Creates a new user story within a specific project.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `201 Created`

### Get All User Stories for a Project

-   **Endpoint:** `GET /api/projects/:id/userstories`
-   **Description:** Retrieves a list of all user stories for a specific project (the Product Backlog).
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK`

### Get User Story by ID

-   **Endpoint:** `GET /api/userstories/:storyId`
-   **Description:** Retrieves the full details of a single user story, including its related Project and Sprint.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK`

### Update User Story

-   **Endpoint:** `PUT /api/userstories/:storyId`
-   **Description:** Updates the details of an existing user story.
-   **Access:** Authenticated (Platform Admin, or Project's `product_owner` / `scrum_master`)
-   **Success Response:** `200 OK`

### Delete User Story

-   **Endpoint:** `DELETE /api/userstories/:storyId`
-   **Description:** Deletes a user story from the product backlog.
-   **Access:** Authenticated (Platform Admin, or Project's `product_owner` / `scrum_master`)
-   **Success Response:** `204 No Content`

---

## 5. Sprints

### Create Sprint

-   **Endpoint:** `POST /api/projects/:id/sprints`
-   **Description:** Creates a new sprint within a specific project.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `201 Created`

### Get All Sprints for a Project

-   **Endpoint:** `GET /api/projects/:id/sprints`
-   **Description:** Retrieves a list of all sprints for a specific project.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK`

### Assign User Story to Sprint

-   **Endpoint:** `POST /api/sprints/:sprintId/userstories`
-   **Description:** Assigns an existing user story to a sprint, creating the Sprint Backlog.
-   **Access:** Authenticated (Platform Admin, or Project's `product_owner` / `scrum_master`)
-   **Request Body:**
    ```json
    {
      "userStoryId": 1
    }
    ```
-   **Success Response:** `200 OK`

---

## 6. Tasks

### Create Task

-   **Endpoint:** `POST /api/userstories/:storyId/tasks`
-   **Description:** Creates a new task for a specific user story.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `201 Created`

### Get All Tasks for a User Story

-   **Endpoint:** `GET /api/userstories/:storyId/tasks`
-   **Description:** Retrieves a list of all tasks for a specific user story.
-   **Access:** Authenticated (any valid user)
-   **Success Response:** `200 OK`

### Assign Task to User

-   **Endpoint:** `PUT /api/tasks/:taskId/assign`
-   **Description:** Assigns a task to a user who is a member of the project.
-   **Access:** Authenticated (any valid user)
-   **Request Body:**
    ```json
    {
      "userId": 2
    }
    ```
-   **Success Response:** `200 OK`

### Update Task Status

-   **Endpoint:** `PUT /api/tasks/:taskId/status`
-   **Description:** Updates the status of a task.
-   **Access:** Authenticated (any valid user)
-   **Request Body:**
    ```json
    {
      "status": "in_progress"
    }
    ```
    *Valid statuses are: `todo`, `in_progress`, `in_review`, `done`.*
-   **Success Response:** `200 OK`

---

## 7. Administration (Admin-Only)

All endpoints in this section require the user to have an `admin` platform role.

### Get All Users

-   **Endpoint:** `GET /api/admin/users`
-   **Description:** Retrieves a list of all users.
-   **Access:** Admin only
-   **Success Response:** `200 OK`

### Create Standard User

-   **Endpoint:** `POST /api/admin/users`
-   **Description:** Creates a new standard platform user.
-   **Access:** Admin only
-   **Success Response:** `201 Created`

### Create Admin User

-   **Endpoint:** `POST /api/admin/users/admin`
-   **Description:** Creates a new administrator user.
-   **Access:** Admin only
-   **Success Response:** `201 Created`

### Add Member to Project

-   **Endpoint:** `POST /api/admin/projects/:id/members`
-   **Description:** Assigns a user to a project with a specific role.
-   **Access:** Admin only
-   **Request Body:**
    ```json
    {
      "userId": 4,
      "role": "team_developer"
    }
    ```
    *Valid roles are: `scrum_master`, `product_owner`, `team_developer`.*
-   **Success Response:** `201 Created`
