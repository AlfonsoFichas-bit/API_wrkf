# Gemini API Documentation

## Base URL

The base URL for all API endpoints is:

```
http://localhost:8080
```

This document provides a detailed description of the API endpoints, including request and response formats, and information about middleware.

## API Documentation

### Swagger UI

#### GET /swagger/*

Access the interactive API documentation using Swagger UI.

-   **Description:** Provides an interactive web interface for exploring and testing the API endpoints.
-   **Access URL:** `http://localhost:8080/swagger/index.html`

## Authentication and Middleware

### JWT Authentication

All endpoints under the `/api` prefix require a valid JSON Web Token (JWT) to be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <your_jwt_token>
```

The JWT is obtained by authenticating via the `POST /login` endpoint. The token contains user information, including the user ID and role, which are used for authorization purposes throughout the API.

### Admin Authorization

Endpoints under the `/api/admin` prefix are further protected by an admin-only authorization middleware. This means that only users with the `admin` role can access these resources.

## Endpoints

### Authentication

#### POST /login

Authenticates a user and returns a JWT token.

-   **Request Body:**

    ```json
    {
        "correo": "user@example.com",
        "contraseña": "password123"
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "token": "your_jwt_token"
    }
    ```

#### POST /api/logout

Logs out the user. This endpoint is used to signal the server that the user is logging out. The client is responsible for deleting the JWT token from storage.

-   **Response (200 OK):**

    ```json
    {
        "message": "Logged out successfully"
    }
    ```

### Users

#### GET /api/me

Retrieves the details of the currently authenticated user.

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Nombre": "John",
        "ApellidoPaterno": "Doe",
        "ApellidoMaterno": "Smith",
        "Correo": "user@example.com",
        "Role": "user",
        "CreatedAt": "2023-10-27T10:00:00Z"
    }
    ```

#### GET /api/users/:id

Retrieves a user by their ID.

-   **Path Parameters:**
    -   `id` (uint): The ID of the user to retrieve.
-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Nombre": "John",
        "ApellidoPaterno": "Doe",
        "ApellidoMaterno": "Smith",
        "Correo": "user@example.com",
        "Role": "user",
        "CreatedAt": "2023-10-27T10:00:00Z"
    }
    ```

### Admin - User Management

#### GET /api/admin/users

Retrieves a list of all users. (Admin only)

- **Response (200 OK):**

```json
[
    {
        "ID": 1,
        "Nombre": "Admin",
        "ApellidoPaterno": "User",
        "ApellidoMaterno": "",
        "Correo": "admin@example.com",
        "Role": "admin",
        "CreatedAt": "2023-10-27T10:10:00Z"
    },
    {
        "ID": 2,
        "Nombre": "Jane",
        "ApellidoPaterno": "Doe",
        "ApellidoMaterno": "Smith",
        "Correo": "jane.doe@example.com",
        "Role": "user",
        "CreatedAt": "2023-10-27T10:05:00Z"
    }
]
```

#### POST /api/admin/users

Creates a new standard platform user. (Admin only)

-   **Request Body:**

    ```json
    {
        "Nombre": "Jane",
        "ApellidoPaterno": "Doe",
        "ApellidoMaterno": "Smith",
        "Correo": "jane.doe@example.com",
        "Contraseña": "password123"
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 2,
        "Nombre": "Jane",
        "ApellidoPaterno": "Doe",
        "ApellidoMaterno": "Smith",
        "Correo": "jane.doe@example.com",
        "Role": "user",
        "CreatedAt": "2023-10-27T10:05:00Z"
    }
    ```

#### POST /api/admin/users/admin

Creates a new admin user. (Admin only)

-   **Request Body:**

    ```json
    {
        "Nombre": "Admin",
        "ApellidoPaterno": "User",
        "ApellidoMaterno": "",
        "Correo": "admin@example.com",
        "Contraseña": "adminpassword"
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 3,
        "Nombre": "Admin",
        "ApellidoPaterno": "User",
        "ApellidoMaterno": "",
        "Correo": "admin@example.com",
        "Role": "admin",
        "CreatedAt": "2023-10-27T10:10:00Z"
    }
    ```

#### PUT /api/admin/users/:id

Updates a user's information. (Admin only)

-   **Path Parameters:**
    -   `id` (uint): The ID of the user to update.
-   **Request Body:**

    ```json
    {
        "Nombre": "Updated Name"
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 2,
        "Nombre": "Updated Name",
        "ApellidoPaterno": "Doe",
        "ApellidoMaterno": "Smith",
        "Correo": "jane.doe@example.com",
        "Role": "user",
        "CreatedAt": "2023-10-27T10:05:00Z"
    }
    ```

#### DELETE /api/admin/users/:id

Deletes a user. (Admin only)

-   **Path Parameters:**
    -   `id` (uint): The ID of the user to delete.
-   **Response (204 No Content)**

### Public Admin Creation

#### POST /create-admin

Creates a new admin user. This is a public endpoint and should be used with caution.

-   **Request Body:**

    ```json
    {
        "Nombre": "Admin",
        "ApellidoPaterno": "User",
        "ApellidoMaterno": "",
        "Correo": "newadmin@example.com",
        "Contraseña": "strongpassword"
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 4,
        "Nombre": "Admin",
        "ApellidoPaterno": "User",
        "ApellidoMaterno": "",
        "Correo": "newadmin@example.com",
        "Role": "admin",
        "CreatedAt": "2023-10-27T10:12:00Z"
    }
    ```

### WebSocket Real-time Updates

#### GET /ws

Establishes a WebSocket connection for real-time updates. Ideal for keeping Kanban boards and other interfaces synchronized.

-   **Authentication:** JWT token required as query parameter.
-   **Connection URL:** `ws://localhost:8080/ws?token=YOUR_JWT_TOKEN`
-   **Protocol:** WebSocket (ws://) or Secure WebSocket (wss://) in production.

##### Message Format

All messages follow this JSON structure:

```json
{
  "type": "event_name",
  "payload": { ... }
}
```

##### Supported Events

###### `task_created`

Emitted when a new task is created.

**Payload:**
```json
{
  "task": {
    "ID": 1,
    "Title": "New Task",
    "Description": "Task description",
    "UserStoryID": 1,
    "Status": "todo",
    "AssignedToID": null,
    "EstimatedHours": null,
    "SpentHours": null,
    "IsDeliverable": false,
    "CreatedAt": "2023-10-27T10:55:00Z",
    "UpdatedAt": "2023-10-27T10:55:00Z"
  },
  "timestamp": "2023-10-27T10:55:00Z"
}
```

###### `task_status_updated`

Emitted when a task's status changes.

**Payload:**
```json
{
  "taskId": 1,
  "oldStatus": "todo",
  "newStatus": "in_progress",
  "updatedBy": {
    "id": 2,
    "name": "John Doe"
  },
  "timestamp": "2023-10-27T11:00:00Z"
}
```

###### `task_assigned`

Emitted when a task is assigned to a user.

**Payload:**
```json
{
  "taskId": 1,
  "assignedTo": {
    "id": 3,
    "name": "Jane Smith"
  },
  "assignedBy": {
    "id": 2,
    "name": "John Doe"
  },
  "timestamp": "2023-10-27T11:05:00Z"
}
```

###### `task_deleted`

Emitted when a task is deleted.

**Payload:**
```json
{
  "taskId": 1,
  "deletedBy": {
    "id": 2,
    "name": "John Doe"
  },
  "timestamp": "2023-10-27T11:10:00Z"
}
```

### Projects

#### POST /api/projects

Creates a new project.

-   **Request Body:**

    ```json
    {
        "Name": "New Project",
        "Description": "This is a new project."
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 1,
        "Name": "New Project",
        "Description": "This is a new project.",
        "Status": "planning",
        "StartDate": null,
        "EndDate": null,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:15:00Z",
        "UpdatedAt": "2023-10-27T10:15:00Z"
    }
    ```

#### GET /api/projects

Retrieves all projects.

-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 1,
            "Name": "New Project",
            "Description": "This is a new project.",
            "Status": "planning",
            "StartDate": null,
            "EndDate": null,
            "CreatedByID": 1,
            "CreatedBy": {
                "ID": 1,
                "Nombre": "John",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "john@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:00:00Z"
            },
            "CreatedAt": "2023-10-27T10:15:00Z",
            "UpdatedAt": "2023-10-27T10:15:00Z"
        }
    ]
    ```

#### GET /api/projects/:id

Retrieves a project by its ID.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project to retrieve.
-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Name": "New Project",
        "Description": "This is a new project.",
        "Status": "planning",
        "StartDate": null,
        "EndDate": null,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "Members": [
            {
                "ID": 1,
                "UserID": 2,
                "User": {
                    "ID": 2,
                    "Nombre": "Jane",
                    "ApellidoPaterno": "Doe",
                    "ApellidoMaterno": "Smith",
                    "Correo": "jane@example.com",
                    "Role": "user",
                    "CreatedAt": "2023-10-27T10:05:00Z"
                },
                "ProjectID": 1,
                "Role": "team_developer",
                "CreatedAt": "2023-10-27T10:25:00Z",
                "UpdatedAt": "2023-10-27T10:25:00Z"
            }
        ],
        "CreatedAt": "2023-10-27T10:15:00Z",
        "UpdatedAt": "2023-10-27T10:15:00Z"
    }
    ```

#### PUT /api/projects/:id

Updates a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project to update.
-   **Request Body:**

    ```json
    {
        "Name": "Updated Project Name"
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Name": "Updated Project Name",
        "Description": "This is a new project.",
        "Status": "planning",
        "StartDate": null,
        "EndDate": null,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:15:00Z",
        "UpdatedAt": "2023-10-27T10:20:00Z"
    }
    ```

#### DELETE /api/projects/:id

Deletes a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project to delete.
-   **Response (204 No Content)**

#### GET /api/projects/:id/active-sprint

Retrieves the currently active sprint for a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Name": "Sprint 1",
        "Goal": "Complete feature X",
        "ProjectID": 1,
        "Status": "active",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z",
        "CreatedByID": 1,
        "CreatedAt": "2023-10-27T10:30:00Z",
        "UpdatedAt": "2023-10-27T10:30:00Z"
    }
    ```

-   **Response (404 Not Found):**

    ```json
    {
        "error": "no active sprint found for this project"
    }
    ```

#### GET /api/projects/:id/unassigned-users

Retrieves a list of users who can be added to a project. This list excludes existing project members and any users with the 'admin' role.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 3,
            "Nombre": "Carlos",
            "ApellidoPaterno": "Ruiz",
            "ApellidoMaterno": "Gomez",
            "Correo": "carlos.ruiz@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:11:00Z"
        }
    ]
    ```


### Projects

#### GET /api/projects/:id/members

Retrieves a list of all members for a specific project, including their user details and role within the project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 1,
            "UserID": 2,
            "User": {
                "ID": 2,
                "Nombre": "Jane",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "jane.doe@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:05:00Z"
            },
            "ProjectID": 1,
            "Role": "team_developer",
            "CreatedAt": "2023-10-27T10:25:00Z",
            "UpdatedAt": "2023-10-27T10:25:00Z"
        }
    ]
    ```

### Admin - Project Management

#### POST /api/admin/projects/:id/members

Adds a member to a project. (Admin only)

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Request Body:**

    ```json
    {
        "userId": 2,
        "role": "team_developer"
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 1,
        "UserID": 2,
        "ProjectID": 1,
        "Role": "team_developer",
        "CreatedAt": "2023-10-27T10:25:00Z",
        "UpdatedAt": "2023-10-27T10:25:00Z"
    }
    ```

-   **Response (409 Conflict):**

    Returned if the user is already a member of the project.

    ```json
    {
        "error": "user is already a member of this project"
    }
    ```

### Sprints

#### POST /api/projects/:id/sprints

Creates a new sprint within a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Request Body:**

    El campo `status` es opcional. Si no se incluye, el sprint se creará con el estado `planned` por defecto.

    ```json
    {
        "Name": "Sprint 1",
        "Goal": "Complete feature X",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z"
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 1,
        "Name": "Sprint 1",
        "Goal": "Complete user authentication feature",
        "ProjectID": 1,
        "Project": {
            "ID": 1,
            "Name": "Project Name",
            "Description": "Project description",
            "Status": "active"
        },
        "Status": "planned",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z",
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:30:00Z",
        "UpdatedAt": "2023-10-27T10:30:00Z"
    }
    ```

#### GET /api/projects/:id/sprints

Retrieves all sprints for a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 1,
            "Name": "Sprint 1",
            "Goal": "Complete user authentication feature",
            "ProjectID": 1,
            "Status": "active",
            "StartDate": "2023-11-01T00:00:00Z",
            "EndDate": "2023-11-15T23:59:59Z",
            "CreatedByID": 1,
            "CreatedBy": {
                "ID": 1,
                "Nombre": "John",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "john@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:00:00Z"
            },
            "CreatedAt": "2023-10-27T10:30:00Z",
            "UpdatedAt": "2023-10-27T10:30:00Z"
        }
    ]
    ```

#### GET /api/sprints/:sprintId

Retrieves a sprint by its ID.

-   **Path Parameters:**
    -   `sprintId` (uint): The ID of the sprint.
-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Name": "Sprint 1",
        "Goal": "Complete user authentication feature",
        "ProjectID": 1,
        "Project": {
            "ID": 1,
            "Name": "Project Name",
            "Description": "Project description",
            "Status": "active"
        },
        "Status": "active",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z",
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:30:00Z",
        "UpdatedAt": "2023-10-27T10:30:00Z"
    }
    ```

#### PUT /api/sprints/:sprintId

Updates a sprint.

-   **Path Parameters:**
    -   `sprintId` (uint): The ID of the sprint.
-   **Request Body:**

    ```json
    {
        "Name": "Updated Sprint Name"
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Name": "Updated Sprint Name",
        "Goal": "Complete feature X",
        "ProjectID": 1,
        "Status": "planned",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z",
        "CreatedByID": 1,
        "CreatedAt": "2023-10-27T10:30:00Z",
        "UpdatedAt": "2023-10-27T10:35:00Z"
    }
    ```

#### DELETE /api/sprints/:sprintId

Deletes a sprint.

-   **Path Parameters:**
    -   `sprintId` (uint): The ID of the sprint.
-   **Response (204 No Content)**

#### GET /api/sprints/:sprintId/tasks

Retrieves all tasks for a specific sprint with full relationships (user stories, assigned users, etc.).

-   **Path Parameters:**
    -   `sprintId` (uint): The ID of the sprint.
-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 1,
            "Title": "Implement user authentication",
            "Description": "Create login and registration functionality",
            "UserStoryID": 1,
            "UserStory": {
                "ID": 1,
                "Title": "User Authentication",
                "Description": "As a user, I want to authenticate",
                "AcceptanceCriteria": "Given..., When..., Then...",
                "Priority": "high",
                "Status": "in_sprint",
                "Points": 5,
                "ProjectID": 1,
                "SprintID": 1
            },
            "Status": "in_progress",
            "AssignedToID": 2,
            "AssignedTo": {
                "ID": 2,
                "Nombre": "Jane",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "jane@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:05:00Z"
            },
            "EstimatedHours": 8.5,
            "SpentHours": 3.5,
            "IsDeliverable": true,
            "CreatedByID": 1,
            "CreatedBy": {
                "ID": 1,
                "Nombre": "John",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "john@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:00:00Z"
            },
            "CreatedAt": "2023-10-27T10:55:00Z",
            "UpdatedAt": "2023-10-27T11:00:00Z"
        }
    ]
    ```

#### PUT /api/sprints/:sprintId/status

Updates the status of a sprint (planned/active/completed/cancelled).

-   **Path Parameters:**
    -   `sprintId` (uint): The ID of the sprint.
-   **Request Body:**

    ```json
    {
        "status": "active"
    }
    ```

    Valid status values: `planned`, `active`, `completed`, `cancelled`

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Name": "Sprint 1",
        "Goal": "Complete feature X",
        "ProjectID": 1,
        "Status": "active",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z",
        "CreatedByID": 1,
        "CreatedAt": "2023-10-27T10:30:00Z",
        "UpdatedAt": "2023-10-27T10:35:00Z"
    }
    ```

#### POST /api/sprints/:sprintId/userstories

Assigns a user story to a sprint.

-   **Path Parameters:**
    -   `sprintId` (uint): The ID of the sprint.
-   **Request Body:**

    ```json
    {
        "userStoryId": 1
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "User Story 1",
        "Description": "As a user, I want to...",
        "AcceptanceCriteria": "...",
        "Priority": "high",
        "Status": "in_sprint",
        "Points": 5,
        "ProjectID": 1,
        "SprintID": 1,
        "CreatedByID": 1,
        "CreatedAt": "2023-10-27T10:40:00Z",
        "UpdatedAt": "2023-10-27T10:45:00Z"
    }
    ```

### User Stories

#### POST /api/projects/:id/userstories

Creates a new user story within a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Request Body:**

    ```json
    {
        "Title": "User Story 1",
        "Description": "As a user, I want to...",
        "AcceptanceCriteria": "..."
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 1,
        "Title": "User Story 1",
        "Description": "As a user, I want to...",
        "AcceptanceCriteria": "Given..., When..., Then...",
        "Priority": "high",
        "Status": "backlog",
        "Points": 5,
        "ProjectID": 1,
        "Project": {
            "ID": 1,
            "Name": "Project Name",
            "Description": "Project description",
            "Status": "planning"
        },
        "SprintID": null,
        "Sprint": null,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "AssignedToID": null,
        "AssignedTo": null,
        "CreatedAt": "2023-10-27T10:40:00Z",
        "UpdatedAt": "2023-10-27T10:40:00Z"
    }
    ```

#### GET /api/projects/:id/userstories

Retrieves all user stories for a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 1,
            "Title": "User Story 1",
            "Description": "As a user, I want to...",
            "AcceptanceCriteria": "Given..., When..., Then...",
            "Priority": "high",
            "Status": "backlog",
            "Points": 5,
            "ProjectID": 1,
            "SprintID": null,
            "CreatedByID": 1,
            "CreatedBy": {
                "ID": 1,
                "Nombre": "John",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "john@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:00:00Z"
            },
            "CreatedAt": "2023-10-27T10:40:00Z",
            "UpdatedAt": "2023-10-27T10:40:00Z"
        }
    ]
    ```

#### GET /api/userstories/:storyId

Retrieves a user story by its ID.

-   **Path Parameters:**
    -   `storyId` (uint): The ID of the user story.
-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "User Story 1",
        "Description": "As a user, I want to...",
        "AcceptanceCriteria": "Given..., When..., Then...",
        "Priority": "high",
        "Status": "in_sprint",
        "Points": 5,
        "ProjectID": 1,
        "Project": {
            "ID": 1,
            "Name": "Project Name",
            "Description": "Project description",
            "Status": "active"
        },
        "SprintID": 1,
        "Sprint": {
            "ID": 1,
            "Name": "Sprint 1",
            "Goal": "Complete user authentication",
            "ProjectID": 1,
            "Status": "active",
            "StartDate": "2023-11-01T00:00:00Z",
            "EndDate": "2023-11-15T23:59:59Z"
        },
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "AssignedToID": 2,
        "AssignedTo": {
            "ID": 2,
            "Nombre": "Jane",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "jane@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:05:00Z"
        },
        "CreatedAt": "2023-10-27T10:40:00Z",
        "UpdatedAt": "2023-10-27T10:45:00Z"
    }
    ```

#### PUT /api/userstories/:storyId

Updates a user story.

-   **Path Parameters:**
    -   `storyId` (uint): The ID of the user story.
-   **Request Body:** (partial update supported)

    ```json
    {
        "Title": "Updated User Story Title",
        "Description": "Updated description",
        "AcceptanceCriteria": "Updated acceptance criteria",
        "Priority": "high",
        "Points": 8
    }
    ```

    **Updatable fields:** `Title`, `Description`, `AcceptanceCriteria`, `Priority`, `Points`

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "Updated User Story Title",
        "Description": "Updated description",
        "AcceptanceCriteria": "Updated acceptance criteria",
        "Priority": "high",
        "Status": "in_sprint",
        "Points": 8,
        "ProjectID": 1,
        "SprintID": 1,
        "CreatedByID": 1,
        "CreatedAt": "2023-10-27T10:40:00Z",
        "UpdatedAt": "2023-10-27T10:50:00Z"
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "Updated User Story Title",
        "Description": "As a user, I want to...",
        "AcceptanceCriteria": "Given..., When..., Then...",
        "Priority": "high",
        "Status": "in_sprint",
        "Points": 5,
        "ProjectID": 1,
        "Project": {
            "ID": 1,
            "Name": "Project Name",
            "Description": "Project description",
            "Status": "active"
        },
        "SprintID": 1,
        "Sprint": {
            "ID": 1,
            "Name": "Sprint 1",
            "Goal": "Complete user authentication",
            "ProjectID": 1,
            "Status": "active",
            "StartDate": "2023-11-01T00:00:00Z",
            "EndDate": "2023-11-15T23:59:59Z"
        },
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "AssignedToID": 2,
        "AssignedTo": {
            "ID": 2,
            "Nombre": "Jane",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "jane@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:05:00Z"
        },
        "CreatedAt": "2023-10-27T10:40:00Z",
        "UpdatedAt": "2023-10-27T10:50:00Z"
    }
    ```

#### DELETE /api/userstories/:storyId

Deletes a user story.

-   **Path Parameters:**
    -   `storyId` (uint): The ID of the user story.
-   **Response (204 No Content)**

### Tasks

#### POST /api/userstories/:storyId/tasks

Creates a new task within a user story.

-   **Path Parameters:**
    -   `storyId` (uint): The ID of the user story.
-   **Request Body:**

    ```json
    {
        "Title": "Task 1",
        "Description": "Implement the first part of the story."
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 1,
        "Title": "Implement user authentication",
        "Description": "Create login and registration functionality",
        "UserStoryID": 1,
        "UserStory": {
            "ID": 1,
            "Title": "User Story 1",
            "Description": "As a user, I want to...",
            "Priority": "high",
            "Status": "in_sprint",
            "Points": 5
        },
        "Status": "todo",
        "AssignedToID": null,
        "AssignedTo": null,
        "EstimatedHours": 8.5,
        "SpentHours": 0,
        "IsDeliverable": true,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:55:00Z",
        "UpdatedAt": "2023-10-27T10:55:00Z"
    }
    ```

#### GET /api/userstories/:storyId/tasks

Retrieves all tasks for a user story.

-   **Path Parameters:**
    -   `storyId` (uint): The ID of the user story.
-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 1,
            "Title": "Implement user authentication",
            "Description": "Create login and registration functionality",
            "UserStoryID": 1,
            "UserStory": {
                "ID": 1,
                "Title": "User Story 1",
                "Description": "As a user, I want to...",
                "Priority": "high",
                "Status": "in_sprint",
                "Points": 5
            },
            "Status": "in_progress",
            "AssignedToID": 2,
            "AssignedTo": {
                "ID": 2,
                "Nombre": "Jane",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "jane@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:05:00Z"
            },
            "EstimatedHours": 8.5,
            "SpentHours": 2.5,
            "IsDeliverable": true,
            "CreatedByID": 1,
            "CreatedBy": {
                "ID": 1,
                "Nombre": "John",
                "ApellidoPaterno": "Doe",
                "ApellidoMaterno": "Smith",
                "Correo": "john@example.com",
                "Role": "user",
                "CreatedAt": "2023-10-27T10:00:00Z"
            },
            "CreatedAt": "2023-10-27T10:55:00Z",
            "UpdatedAt": "2023-10-27T11:00:00Z"
        }
    ]
    ```

#### PUT /api/tasks/:taskId

Updates a task.

-   **Path Parameters:**
    -   `taskId` (uint): The ID of the task.
-   **Request Body:**

    ```json
    {
        "Title": "Updated Task Title"
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "Updated Task Title",
        "Description": "Updated description",
        "UserStoryID": 1,
        "UserStory": {
            "ID": 1,
            "Title": "User Story 1",
            "Description": "As a user, I want to...",
            "Priority": "high",
            "Status": "in_sprint",
            "Points": 5
        },
        "Status": "in_progress",
        "AssignedToID": 2,
        "AssignedTo": {
            "ID": 2,
            "Nombre": "Jane",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "jane@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:05:00Z"
        },
        "EstimatedHours": 12,
        "SpentHours": 4,
        "IsDeliverable": true,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:55:00Z",
        "UpdatedAt": "2023-10-27T11:00:00Z"
    }
    ```

#### DELETE /api/tasks/:taskId

Deletes a task.

-   **Path Parameters:**
    -   `taskId` (uint): The ID of the task.
-   **Response (204 No Content)**

#### PUT /api/tasks/:taskId/assign

Assigns a task to a user.

-   **Path Parameters:**
    -   `taskId` (uint): The ID of the task.
-   **Request Body:**

    ```json
    {
        "userId": 2
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "Implement user authentication",
        "Description": "Create login and registration functionality",
        "UserStoryID": 1,
        "UserStory": {
            "ID": 1,
            "Title": "User Story 1",
            "Description": "As a user, I want to...",
            "Priority": "high",
            "Status": "in_sprint",
            "Points": 5
        },
        "Status": "in_progress",
        "AssignedToID": 2,
        "AssignedTo": {
            "ID": 2,
            "Nombre": "Jane",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "jane@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:05:00Z"
        },
        "EstimatedHours": 8.5,
        "SpentHours": 2.5,
        "IsDeliverable": true,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:55:00Z",
        "UpdatedAt": "2023-10-27T11:05:00Z"
    }
    ```

#### PUT /api/tasks/:taskId/status

Updates the status of a task.

-   **Path Parameters:**
    -   `taskId` (uint): The ID of the task.
-   **Request Body:**

    ```json
    {
        "status": "in_progress"
    }
    ```

    Valid status values: `todo`, `in_progress`, `in_review`, `done`

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "Implement user authentication",
        "Description": "Create login and registration functionality",
        "UserStoryID": 1,
        "UserStory": {
            "ID": 1,
            "Title": "User Story 1",
            "Description": "As a user, I want to...",
            "Priority": "high",
            "Status": "in_sprint",
            "Points": 5
        },
        "Status": "in_progress",
        "AssignedToID": 2,
        "AssignedTo": {
            "ID": 2,
            "Nombre": "Jane",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "jane@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:05:00Z"
        },
        "EstimatedHours": 8.5,
        "SpentHours": 3.5,
        "IsDeliverable": true,
        "CreatedByID": 1,
        "CreatedBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "CreatedAt": "2023-10-27T10:55:00Z",
        "UpdatedAt": "2023-10-27T11:10:00Z"
    }
    ```

#### POST /api/tasks/:id/comments

Adds a comment to a task.

-   **Path Parameters:**
    -   `id` (uint): The ID of the task.
-   **Request Body:**

    ```json
    {
        "content": "This is a comment."
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "ID": 1,
        "TaskID": 1,
        "AuthorID": 1,
        "Author": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "Content": "This is a comment about the task implementation.",
        "CreatedAt": "2023-10-27T11:15:00Z",
        "UpdatedAt": "2023-10-27T11:15:00Z"
    }
    ```

### Notifications

#### GET /api/notifications

Retrieves all notifications for the current user.

-   **Response (200 OK):**

    ```json
    [
        {
            "ID": 1,
            "user_id": 1,
            "message": "You have been assigned to a new task.",
            "is_read": false,
            "link": "/tasks/1"
        }
    ]
    ```

#### POST /api/notifications/read/all

Marks all of a user's notifications as read.

-   **Response (204 No Content)**

#### POST /api/notifications/:id/read

Marks a notification as read.

-   **Path Parameters:**
    -   `id` (uint): The ID of the notification.
-   **Response (204 No Content)**

### Rubrics

#### POST /api/rubrics

Creates a new rubric.

-   **Request Body:**

    ```json
    {
        "name": "New Rubric",
        "description": "This is a new rubric.",
        "projectId": 1,
        "isTemplate": false,
        "criteria": [
            {
                "title": "Criterion 1",
                "description": "Description for criterion 1",
                "maxPoints": 10,
                "levels": [
                    {
                        "score": 10,
                        "description": "Excellent"
                    },
                    {
                        "score": 5,
                        "description": "Average"
                    }
                ]
            }
        ]
    }
    ```

-   **Response (201 Created):**

    ```json
    {
        "id": 1,
        "name": "Code Quality Rubric",
        "description": "Evaluates code quality standards",
        "projectId": 1,
        "project": {
            "ID": 1,
            "Name": "Project Name",
            "Description": "Project description",
            "Status": "active"
        },
        "createdById": 1,
        "createdBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "status": "DRAFT",
        "isTemplate": false,
        "criteria": [
            {
                "id": 1,
                "rubricId": 1,
                "title": "Code Structure",
                "description": "Organization and structure of the code",
                "maxPoints": 10,
                "levels": [
                    {
                        "id": 1,
                        "criterionId": 1,
                        "score": 10,
                        "description": "Excellent structure and organization"
                    },
                    {
                        "id": 2,
                        "criterionId": 1,
                        "score": 7,
                        "description": "Good structure with minor issues"
                    },
                    {
                        "id": 3,
                        "criterionId": 1,
                        "score": 4,
                        "description": "Poor structure, needs improvement"
                    }
                ],
                "createdAt": "2023-10-27T12:00:00Z",
                "updatedAt": "2023-10-27T12:00:00Z"
            }
        ],
        "createdAt": "2023-10-27T12:00:00Z",
        "updatedAt": "2023-10-27T12:00:00Z"
    }
    ```

#### GET /api/rubrics

Retrieves all rubrics.

-   **Query Parameters (optional):**
    -   `isTemplate` (boolean): Filter by template status.
-   **Response (200 OK):**

    ```json
    [
        {
            "id": 1,
            "name": "New Rubric",
            "description": "This is a new rubric.",
            "projectId": 1,
            "createdById": 1,
            "status": "DRAFT",
            "isTemplate": false
        }
    ]
    ```

#### GET /api/rubrics/:id

Retrieves a rubric by its ID.

-   **Path Parameters:**
    -   `id` (uint): The ID of the rubric.
-   **Response (200 OK):**

    ```json
    {
        "id": 1,
        "name": "Code Quality Rubric",
        "description": "Evaluates code quality standards",
        "projectId": 1,
        "project": {
            "ID": 1,
            "Name": "Project Name",
            "Description": "Project description",
            "Status": "active"
        },
        "createdById": 1,
        "createdBy": {
            "ID": 1,
            "Nombre": "John",
            "ApellidoPaterno": "Doe",
            "ApellidoMaterno": "Smith",
            "Correo": "john@example.com",
            "Role": "user",
            "CreatedAt": "2023-10-27T10:00:00Z"
        },
        "status": "DRAFT",
        "isTemplate": false,
        "criteria": [
            {
                "id": 1,
                "rubricId": 1,
                "title": "Code Structure",
                "description": "Organization and structure of the code",
                "maxPoints": 10,
                "levels": [
                    {
                        "id": 1,
                        "criterionId": 1,
                        "score": 10,
                        "description": "Excellent structure and organization"
                    },
                    {
                        "id": 2,
                        "criterionId": 1,
                        "score": 7,
                        "description": "Good structure with minor issues"
                    },
                    {
                        "id": 3,
                        "criterionId": 1,
                        "score": 4,
                        "description": "Poor structure, needs improvement"
                    }
                ],
                "createdAt": "2023-10-27T12:00:00Z",
                "updatedAt": "2023-10-27T12:00:00Z"
            }
        ],
        "createdAt": "2023-10-27T12:00:00Z",
        "updatedAt": "2023-10-27T12:00:00Z"
    }
    ```

#### PUT /api/rubrics/:id

Updates a rubric.

-   **Path Parameters:**
    -   `id` (uint): The ID of the rubric.
-   **Request Body:**

    Any field of the rubric can be updated.

    ```json
    {
        "name": "Updated Rubric Name",
        "description": "This is an updated description.",
        "isTemplate": true
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "id": 1,
        "name": "Updated Rubric Name",
        "description": "This is a new rubric.",
        "projectId": 1,
        "createdById": 1,
        "status": "DRAFT",
        "isTemplate": false
    }
    ```

#### DELETE /api/rubrics/:id

Deletes a rubric.

-   **Path Parameters:**
    -   `id` (uint): The ID of the rubric.
-   **Response (204 No Content)**

#### POST /api/rubrics/:id/duplicate

Duplicates a rubric.

-   **Path Parameters:**
    -   `id` (uint): The ID of the rubric to duplicate.
-   **Response (201 Created):**

    ```json
    {
        "id": 2,
        "name": "New Rubric (Copy)",
        "description": "This is a new rubric.",
        "projectId": 1,
        "createdById": 1,
        "status": "DRAFT",
        "isTemplate": false
    }
    ```
