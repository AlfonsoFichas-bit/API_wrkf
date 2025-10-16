# Gemini API Documentation

## Base URL

The base URL for all API endpoints is:

```
http://localhost:8080
```

This document provides a detailed description of the API endpoints, including request and response formats, and information about middleware.

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
        "CreatedAt": "2023-10-27T10:15:00Z",
        "UpdatedAt": "2023-10-27T10:20:00Z"
    }
    ```

#### DELETE /api/projects/:id

Deletes a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project to delete.
-   **Response (204 No Content)**

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

### Sprints

#### POST /api/projects/:id/sprints

Creates a new sprint within a project.

-   **Path Parameters:**
    -   `id` (uint): The ID of the project.
-   **Request Body:**

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
        "Goal": "Complete feature X",
        "ProjectID": 1,
        "Status": "planned",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z",
        "CreatedByID": 1,
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
            "Goal": "Complete feature X",
            "ProjectID": 1,
            "Status": "planned",
            "StartDate": "2023-11-01T00:00:00Z",
            "EndDate": "2023-11-15T23:59:59Z",
            "CreatedByID": 1,
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
        "Goal": "Complete feature X",
        "ProjectID": 1,
        "Status": "planned",
        "StartDate": "2023-11-01T00:00:00Z",
        "EndDate": "2023-11-15T23:59:59Z",
        "CreatedByID": 1,
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
        "AcceptanceCriteria": "...",
        "Priority": "medium",
        "Status": "backlog",
        "Points": null,
        "ProjectID": 1,
        "SprintID": null,
        "CreatedByID": 1,
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
            "AcceptanceCriteria": "...",
            "Priority": "medium",
            "Status": "backlog",
            "Points": null,
            "ProjectID": 1,
            "SprintID": null,
            "CreatedByID": 1,
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
        "AcceptanceCriteria": "...",
        "Priority": "medium",
        "Status": "backlog",
        "Points": null,
        "ProjectID": 1,
        "SprintID": null,
        "CreatedByID": 1,
        "CreatedAt": "2023-10-27T10:40:00Z",
        "UpdatedAt": "2023-10-27T10:40:00Z"
    }
    ```

#### PUT /api/userstories/:storyId

Updates a user story.

-   **Path Parameters:**
    -   `storyId` (uint): The ID of the user story.
-   **Request Body:**

    ```json
    {
        "Title": "Updated User Story Title"
    }
    ```

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "Updated User Story Title",
        "Description": "As a user, I want to...",
        "AcceptanceCriteria": "...",
        "Priority": "medium",
        "Status": "backlog",
        "Points": null,
        "ProjectID": 1,
        "SprintID": null,
        "CreatedByID": 1,
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
        "Title": "Task 1",
        "Description": "Implement the first part of the story.",
        "UserStoryID": 1,
        "Status": "todo",
        "AssignedToID": null,
        "EstimatedHours": null,
        "SpentHours": null,
        "IsDeliverable": false,
        "CreatedByID": 1,
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
            "Title": "Task 1",
            "Description": "Implement the first part of the story.",
            "UserStoryID": 1,
            "Status": "todo",
            "AssignedToID": null,
            "EstimatedHours": null,
            "SpentHours": null,
            "IsDeliverable": false,
            "CreatedByID": 1,
            "CreatedAt": "2023-10-27T10:55:00Z",
            "UpdatedAt": "2023-10-27T10:55:00Z"
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
        "Description": "Implement the first part of the story.",
        "UserStoryID": 1,
        "Status": "todo",
        "AssignedToID": null,
        "EstimatedHours": null,
        "SpentHours": null,
        "IsDeliverable": false,
        "CreatedByID": 1,
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
        "Title": "Task 1",
        "Description": "Implement the first part of the story.",
        "UserStoryID": 1,
        "Status": "todo",
        "AssignedToID": 2,
        "EstimatedHours": null,
        "SpentHours": null,
        "IsDeliverable": false,
        "CreatedByID": 1,
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

-   **Response (200 OK):**

    ```json
    {
        "ID": 1,
        "Title": "Task 1",
        "Description": "Implement the first part of the story.",
        "UserStoryID": 1,
        "Status": "in_progress",
        "AssignedToID": 2,
        "EstimatedHours": null,
        "SpentHours": null,
        "IsDeliverable": false,
        "CreatedByID": 1,
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
        "Content": "This is a comment.",
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
        "name": "New Rubric",
        "description": "This is a new rubric.",
        "projectId": 1,
        "createdById": 1,
        "status": "DRAFT",
        "isTemplate": false,
        "criteria": [
            {
                "id": 1,
                "rubricId": 1,
                "title": "Criterion 1",
                "description": "Description for criterion 1",
                "maxPoints": 10,
                "levels": [
                    {
                        "id": 1,
                        "criterionId": 1,
                        "score": 10,
                        "description": "Excellent"
                    },
                    {
                        "id": 2,
                        "criterionId": 1,
                        "score": 5,
                        "description": "Average"
                    }
                ]
            }
        ]
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
        "name": "New Rubric",
        "description": "This is a new rubric.",
        "projectId": 1,
        "createdById": 1,
        "status": "DRAFT",
        "isTemplate": false,
        "criteria": [
            {
                "id": 1,
                "rubricId": 1,
                "title": "Criterion 1",
                "description": "Description for criterion 1",
                "maxPoints": 10,
                "levels": [
                    {
                        "id": 1,
                        "criterionId": 1,
                        "score": 10,
                        "description": "Excellent"
                    },
                    {
                        "id": 2,
                        "criterionId": 1,
                        "score": 5,
                        "description": "Average"
                    }
                ]
            }
        ]
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
