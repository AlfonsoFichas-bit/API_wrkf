# API Documentation

This document provides detailed documentation for the API endpoints.

## API Documentation Access

### Swagger UI

- **Endpoint:** `GET /swagger/*`
- **Description:** Interactive API documentation and testing interface
- **Access URL:** `http://localhost:8080/swagger/index.html`

## Authentication

### POST /login

- **Description:** Authenticates a user and returns a JWT token.
- **Request Body:**
  ```json
  {
    "correo": "user@example.com",
    "contraseña": "password"
  }
  ```
- **Responses:**
  - `200 OK`: Returns a JWT token.
    ```json
    {
      "token": "your_jwt_token"
    }
    ```
  - `400 Bad Request`: Invalid input.
  - `401 Unauthorized`: Invalid credentials.

### POST /create-admin

- **Description:** Creates a new admin user. This is typically a setup endpoint.
- **Request Body:** `models.User`
- **Responses:**
  - `201 Created`: Returns the created admin user (without password).
  - `400 Bad Request`: Invalid input.
  - `500 Internal Server Error`: Could not create admin user.

## Users

### GET /api/me

- **Authentication:** JWT Token required.
- **Description:** Retrieves the details of the currently authenticated user.
- **Responses:**
   - `200 OK`: Returns the current user's data (`models.User`). Password field is excluded for security.
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
   - `401 Unauthorized`: Invalid token.
   - `404 Not Found`: User not found.

### GET /api/users/:id

- **Authentication:** JWT Token required.
- **Description:** Retrieves a specific user by their ID.
- **URL Parameters:**
   - `id` (integer): The ID of the user.
- **Responses:**
   - `200 OK`: Returns the user data (`models.User`). Password field is excluded for security.
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
   - `400 Bad Request`: Invalid ID.
   - `404 Not Found`: User not found.

## Notifications

### GET /api/notifications

- **Authentication:** JWT Token required.
- **Description:** Gets all notifications for the current user.
- **Responses:**
   - `200 OK`: Returns a list of notifications (`[]models.Notification`).
     ```json
     [
         {
             "ID": 1,
             "user_id": 1,
             "message": "You have been assigned to a new task.",
             "is_read": false,
             "link": "/tasks/1",
             "CreatedAt": "2023-10-27T10:00:00Z",
             "UpdatedAt": "2023-10-27T10:00:00Z",
             "DeletedAt": null
         }
     ]
     ```
   - `500 Internal Server Error`: Failed to retrieve notifications.

### POST /api/notifications/read/all

- **Authentication:** JWT Token required.
- **Description:** Marks all of a user's notifications as read.
- **Responses:**
  - `204 No Content`: Successfully marked all as read.
  - `500 Internal Server Error`: Failed to mark all notifications as read.

### POST /api/notifications/:id/read

- **Authentication:** JWT Token required.
- **Description:** Marks a specific notification as read.
- **URL Parameters:**
  - `id` (integer): The ID of the notification.
- **Responses:**
  - `204 No Content`: Successfully marked as read.
  - `400 Bad Request`: Invalid notification ID.
  - `500 Internal Server Error`: Failed to mark notification as read.

## Projects

### POST /api/projects

- **Authentication:** JWT Token required.
- **Description:** Creates a new project.
- **Request Body:** `models.Project`
  ```json
  {
    "Name": "New Project",
    "Description": "This is a new project.",
    "Status": "planning",
    "StartDate": "2023-11-01T00:00:00Z", // optional
    "EndDate": "2023-12-01T00:00:00Z"    // optional
  }
  ```
- **Responses:**
   - `201 Created`: Returns the created project (`models.Project`).
     ```json
     {
       "ID": 1,
       "Name": "New Project",
       "Description": "This is a new project.",
       "Status": "planning",
       "StartDate": "2023-11-01T00:00:00Z",
       "EndDate": "2023-12-01T00:00:00Z",
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
   - `400 Bad Request`: Invalid input.
   - `500 Internal Server Error`: Could not create project.

### GET /api/projects

- **Authentication:** JWT Token required.
- **Description:** Retrieves all projects.
- **Responses:**
   - `200 OK`: Returns a list of projects (`[]models.Project`).
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
   - `500 Internal Server Error`: Could not retrieve projects.

### GET /api/projects/:id

- **Authentication:** JWT Token required.
- **Description:** Retrieves a single project by its ID.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Responses:**
   - `200 OK`: Returns the project data (`models.Project`).
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
   - `400 Bad Request`: Invalid project ID.
   - `404 Not Found`: Project not found.

### PUT /api/projects/:id

- **Authentication:** JWT Token required.
- **Description:** Updates a project. Only project owner or admin can update.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Request Body:** A map of fields to update (partial update supported).
  ```json
  {
    "Name": "Updated Project Name",
    "Description": "Updated description.",
    "Status": "active",
    "StartDate": "2023-11-01T00:00:00Z",
    "EndDate": "2023-12-01T00:00:00Z"
  }
  ```
- **Responses:**
   - `200 OK`: Returns the updated project (`models.Project`).
     ```json
     {
       "ID": 1,
       "Name": "Updated Project Name",
       "Description": "Updated description.",
       "Status": "active",
       "StartDate": "2023-11-01T00:00:00Z",
       "EndDate": "2023-12-01T00:00:00Z",
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
   - `400 Bad Request`: Invalid project ID or request body.
   - `403 Forbidden`: User does not have permission to update.
   - `500 Internal Server Error`: Could not update project.

### DELETE /api/projects/:id

- **Authentication:** JWT Token required.
- **Description:** Deletes a project.
- **URL Parameters:**
  - `id` (integer): The ID of the project.
- **Responses:**
  - `204 No Content`: Successfully deleted the project.
  - `400 Bad Request`: Invalid project ID.
  - `403 Forbidden`: User does not have permission to delete.
  - `500 Internal Server Error`: Could not delete project.

### GET /api/projects/:id/unassigned-users

- **Authentication:** JWT Token required.
- **Description:** Retrieves users who are not yet assigned to the project. This list excludes existing project members and any users with the 'admin' role.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Responses:**
   - `200 OK`: Returns a list of unassigned users (`[]models.User`). Password fields are excluded for security.
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
   - `400 Bad Request`: Invalid project ID.
   - `500 Internal Server Error`: Could not retrieve unassigned users.

### GET /api/projects/:id/members

- **Authentication:** JWT Token required.
- **Description:** Retrieves all members of a project, including their user details and role within the project.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Responses:**
   - `200 OK`: Returns a list of project members (`[]models.ProjectMember`).
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
   - `400 Bad Request`: Invalid project ID.
   - `500 Internal Server Error`: Could not retrieve project members.

## Rubrics

### POST /api/rubrics

- **Authentication:** JWT Token required.
- **Description:** Creates a new rubric.
- **Request Body:** `models.Rubric`
  ```json
  {
    "name": "Code Quality Rubric",
    "description": "Evaluates code quality standards",
    "projectId": 1,
    "isTemplate": false,
    "criteria": [
      {
        "title": "Code Structure",
        "description": "Organization and structure of the code",
        "maxPoints": 10,
        "levels": [
          {
            "score": 10,
            "description": "Excellent structure and organization"
          },
          {
            "score": 7,
            "description": "Good structure with minor issues"
          },
          {
            "score": 4,
            "description": "Poor structure, needs improvement"
          }
        ]
      }
    ]
  }
  ```
- **Responses:**
   - `201 Created`: Returns the created rubric (`models.Rubric`).
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
   - `400 Bad Request`: Invalid request body.
   - `500 Internal Server Error`: Could not create rubric.

### GET /api/rubrics

- **Authentication:** JWT Token required.
- **Description:** Retrieves all rubrics. Can be filtered by query parameters (e.g., `?isTemplate=true`).
- **Responses:**
  - `200 OK`: Returns a list of rubrics (`[]models.Rubric`).
  - `500 Internal Server Error`: Could not retrieve rubrics.

### GET /api/rubrics/:id

- **Authentication:** JWT Token required.
- **Description:** Retrieves a single rubric by its ID.
- **URL Parameters:**
   - `id` (integer): The ID of the rubric.
- **Responses:**
   - `200 OK`: Returns the rubric data (`models.Rubric`).
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
   - `400 Bad Request`: Invalid ID format.
   - `404 Not Found`: Rubric not found.

### PUT /api/rubrics/:id

- **Authentication:** JWT Token required.
- **Description:** Updates an existing rubric.
- **URL Parameters:**
  - `id` (integer): The ID of the rubric.
- **Request Body:** `models.Rubric`
- **Responses:**
  - `200 OK`: Returns the updated rubric (`models.Rubric`).
  - `400 Bad Request`: Invalid ID format or request body.
  - `500 Internal Server Error`: Could not update rubric.

### DELETE /api/rubrics/:id

- **Authentication:** JWT Token required.
- **Description:** Deletes a rubric.
- **URL Parameters:**
  - `id` (integer): The ID of the rubric.
- **Responses:**
  - `204 No Content`: Successfully deleted the rubric.
  - `400 Bad Request`: Invalid ID format.
  - `500 Internal Server Error`: Could not delete rubric.

### POST /api/rubrics/:id/duplicate

- **Authentication:** JWT Token required.
- **Description:** Duplicates an existing rubric.
- **URL Parameters:**
  - `id` (integer): The ID of the rubric to duplicate.
- **Responses:**
  - `201 Created`: Returns the new, duplicated rubric (`models.Rubric`).
  - `400 Bad Request`: Invalid ID format.
  - `500 Internal Server Error`: Could not duplicate rubric.

## User Stories

### POST /api/projects/:id/userstories

- **Authentication:** JWT Token required.
- **Description:** Creates a new user story within a project.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Request Body:** `models.UserStory`
  ```json
  {
    "Title": "User Story 1",
    "Description": "As a user, I want to...",
    "AcceptanceCriteria": "Given..., When..., Then...",
    "Priority": "high",
    "Points": 5
  }
  ```
- **Responses:**
   - `201 Created`: Returns the created user story (`models.UserStory`).
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
   - `400 Bad Request`: Invalid project ID or input.
   - `500 Internal Server Error`: Could not create user story.

### GET /api/projects/:id/userstories

- **Authentication:** JWT Token required.
- **Description:** Retrieves all user stories for a specific project.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Responses:**
   - `200 OK`: Returns a list of user stories (`[]models.UserStory`).
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
   - `400 Bad Request`: Invalid project ID.
   - `500 Internal Server Error`: Could not retrieve user stories.

### GET /api/userstories/:storyId

- **Authentication:** JWT Token required.
- **Description:** Retrieves a single user story by its ID, including related project, sprint, and user data.
- **URL Parameters:**
   - `storyId` (integer): The ID of the user story.
- **Responses:**
   - `200 OK`: Returns the user story data (`models.UserStory`).
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
   - `400 Bad Request`: Invalid user story ID.
   - `404 Not Found`: User story not found.

### PUT /api/userstories/:storyId

- **Authentication:** JWT Token required.
- **Description:** Updates a user story.
- **URL Parameters:**
  - `storyId` (integer): The ID of the user story.
- **Request Body:** A map of fields to update (partial update supported).
  ```json
  {
    "Title": "Updated Title",
    "Description": "Updated description",
    "AcceptanceCriteria": "Updated acceptance criteria",
    "Priority": "high",
    "Points": 8
  }
  ```
  **Updatable fields:** `Title`, `Description`, `AcceptanceCriteria`, `Priority`, `Points`
- **Responses:**
  - `200 OK`: Returns the updated user story (`models.UserStory`).
  - `400 Bad Request`: Invalid user story ID or request body.
  - `403 Forbidden`: User does not have permission to update.
  - `404 Not Found`: User story not found.
  - `500 Internal Server Error`: Could not update user story.

### DELETE /api/userstories/:storyId

- **Authentication:** JWT Token required.
- **Description:** Deletes a user story.
- **URL Parameters:**
  - `storyId` (integer): The ID of the user story.
- **Responses:**
  - `204 No Content`: Successfully deleted the user story.
  - `400 Bad Request`: Invalid user story ID.
  - `403 Forbidden`: User does not have permission to delete.
  - `404 Not Found`: User story not found.
  - `500 Internal Server Error`: Could not delete user story.

## Tasks

### POST /api/userstories/:storyId/tasks

- **Authentication:** JWT Token required.
- **Description:** Creates a new task for a user story.
- **URL Parameters:**
   - `storyId` (integer): The ID of the user story.
- **Request Body:** `models.Task`
  ```json
  {
    "Title": "Implement user authentication",
    "Description": "Create login and registration functionality",
    "EstimatedHours": 8.5,
    "SpentHours": 0,
    "IsDeliverable": true
  }
  ```
- **Responses:**
   - `201 Created`: Returns the created task (`models.Task`).
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
   - `400 Bad Request`: Invalid user story ID or input.
   - `500 Internal Server Error`: Could not create task.

### GET /api/userstories/:storyId/tasks

- **Authentication:** JWT Token required.
- **Description:** Retrieves all tasks for a user story.
- **URL Parameters:**
   - `storyId` (integer): The ID of the user story.
- **Responses:**
   - `200 OK`: Returns a list of tasks (`[]models.Task`).
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
   - `400 Bad Request`: Invalid user story ID.
   - `500 Internal Server Error`: Could not retrieve tasks.

### PUT /api/tasks/:taskId

- **Authentication:** JWT Token required.
- **Description:** Updates a task.
- **URL Parameters:**
   - `taskId` (integer): The ID of the task.
- **Request Body:** `models.Task` (partial update supported)
  ```json
  {
    "Title": "Updated Task Title",
    "Description": "Updated description",
    "EstimatedHours": 12,
    "SpentHours": 4,
    "IsDeliverable": true
  }
  ```
- **Responses:**
   - `200 OK`: Returns the updated task (`models.Task`).
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
   - `400 Bad Request`: Invalid task ID or input.
   - `404 Not Found`: Task not found.
   - `500 Internal Server Error`: Could not update task.

### DELETE /api/tasks/:taskId

- **Authentication:** JWT Token required.
- **Description:** Deletes a task.
- **URL Parameters:**
  - `taskId` (integer): The ID of the task.
- **Responses:**
  - `204 No Content`: Successfully deleted the task.
  - `400 Bad Request`: Invalid task ID.
  - `404 Not Found`: Task not found or could not be deleted.

### PUT /api/tasks/:taskId/assign

- **Authentication:** JWT Token required.
- **Description:** Assigns a task to a user who is a member of the project.
- **URL Parameters:**
   - `taskId` (integer): The ID of the task.
- **Request Body:**
   ```json
   {
     "userId": 2
   }
   ```
- **Responses:**
   - `200 OK`: Returns the updated task (`models.Task`).
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
   - `400 Bad Request`: Invalid task ID or request body.

### PUT /api/tasks/:taskId/status

- **Authentication:** JWT Token required.
- **Description:** Updates the status of a task.
- **URL Parameters:**
   - `taskId` (integer): The ID of the task.
- **Request Body:**
   ```json
   {
     "status": "in_progress"
   }
   ```
   Valid status values: `todo`, `in_progress`, `in_review`, `done`
- **Responses:**
   - `200 OK`: Returns the updated task (`models.Task`).
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
   - `400 Bad Request`: Invalid task ID or request body.
   - `500 Internal Server Error`: Could not update task status.

### POST /api/tasks/:id/comments

- **Authentication:** JWT Token required.
- **Description:** Adds a comment to a task.
- **URL Parameters:**
   - `id` (integer): The ID of the task.
- **Request Body:**
   ```json
   {
     "content": "This is a comment about the task implementation."
   }
   ```
- **Responses:**
   - `201 Created`: Returns the created comment (`models.TaskComment`).
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
   - `400 Bad Request`: Invalid task ID or empty comment.
   - `500 Internal Server Error`: Could not add comment.

## Sprints

### POST /api/projects/:id/sprints

- **Authentication:** JWT Token required.
- **Description:** Creates a new sprint within a project.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Request Body:** `models.Sprint`. El campo `status` es opcional; si no se provee, el sprint se creará con el estado `'planned'`.
  ```json
  {
    "Name": "Sprint 1",
    "Goal": "Complete user authentication feature",
    "StartDate": "2023-11-01T00:00:00Z",
    "EndDate": "2023-11-15T23:59:59Z",
    "Status": "planned"
  }
  ```
- **Responses:**
   - `201 Created`: Returns the created sprint (`models.Sprint`).
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
   - `400 Bad Request`: Invalid project ID or input.
   - `500 Internal Server Error`: Could not create sprint.

### GET /api/projects/:id/sprints

- **Authentication:** JWT Token required.
- **Description:** Retrieves all sprints for a project.
- **URL Parameters:**
   - `id` (integer): The ID of the project.
- **Responses:**
   - `200 OK`: Returns a list of sprints (`[]models.Sprint`).
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
   - `400 Bad Request`: Invalid project ID.
   - `500 Internal Server Error`: Could not retrieve sprints.

### GET /api/sprints/:sprintId

- **Authentication:** JWT Token required.
- **Description:** Retrieves a single sprint by its ID.
- **URL Parameters:**
   - `sprintId` (integer): The ID of the sprint.
- **Responses:**
   - `200 OK`: Returns the sprint data (`models.Sprint`).
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
   - `400 Bad Request`: Invalid sprint ID.
   - `404 Not Found`: Sprint not found.

### PUT /api/sprints/:sprintId

- **Authentication:** JWT Token required.
- **Description:** Updates a sprint.
- **URL Parameters:**
  - `sprintId` (integer): The ID of the sprint.
- **Request Body:** `models.Sprint`
- **Responses:**
  - `200 OK`: Returns the updated sprint (`models.Sprint`).
  - `400 Bad Request`: Invalid sprint ID or input.
  - `404 Not Found`: Sprint not found.
  - `500 Internal Server Error`: Could not update sprint.

### DELETE /api/sprints/:sprintId

- **Authentication:** JWT Token required.
- **Description:** Deletes a sprint.
- **URL Parameters:**
  - `sprintId` (integer): The ID of the sprint.
- **Responses:**
  - `204 No Content`: Successfully deleted the sprint.
  - `400 Bad Request`: Invalid sprint ID.
  - `404 Not Found`: Sprint not found or could not be deleted.

### POST /api/sprints/:sprintId/userstories

- **Authentication:** JWT Token required.
- **Description:** Assigns a user story to a sprint.
- **URL Parameters:**
  - `sprintId` (integer): The ID of the sprint.
- **Request Body:**
  ```json
  {
    "userStoryId": 456
  }
  ```
- **Responses:**
  - `200 OK`: Returns the updated user story (`models.UserStory`).
  - `400 Bad Request`: Invalid sprint ID or request body.
  - `403 Forbidden`: User does not have permission to assign.
  - `404 Not Found`: Sprint or user story not found.
  - `500 Internal Server Error`: Could not assign user story.

## Admin Routes

The following routes are protected and require admin privileges.

### GET /api/admin/users

- **Authentication:** JWT Token and Admin role required.
- **Description:** Retrieves all users in the system.
- **Responses:**
   - `200 OK`: Returns a list of all users (`[]models.User`). Password fields are excluded for security.
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
   - `500 Internal Server Error`: Could not retrieve users.

### POST /api/admin/users

- **Authentication:** JWT Token and Admin role required.
- **Description:** Creates a new standard user.
- **Request Body:** `models.User`
- **Responses:**
   - `201 Created`: Returns the created user (`models.User`). Password field is excluded for security.
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
   - `400 Bad Request`: Invalid input.
   - `500 Internal Server Error`: Could not create user.

### POST /api/admin/users/admin

- **Authentication:** JWT Token and Admin role required.
- **Description:** Creates a new admin user.
- **Request Body:** `models.User`
- **Responses:**
   - `201 Created`: Returns the created admin user (`models.User`). Password field is excluded for security.
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
   - `400 Bad Request`: Invalid input.
   - `500 Internal Server Error`: Could not create admin user.

### PUT /api/admin/users/:id

- **Authentication:** JWT Token and Admin role required.
- **Description:** Updates a user's information.
- **URL Parameters:**
   - `id` (integer): The ID of the user to update.
- **Request Body:** `models.User`
- **Responses:**
   - `200 OK`: Returns the updated user (`models.User`). Password field is excluded for security.
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
   - `400 Bad Request`: Invalid ID or input.
   - `500 Internal Server Error`: Could not update user.

### DELETE /api/admin/users/:id

- **Authentication:** JWT Token and Admin role required.
- **Description:** Deletes a user.
- **URL Parameters:**
  - `id` (integer): The ID of the user to delete.
- **Responses:**
  - `204 No Content`: Successfully deleted the user.
  - `400 Bad Request`: Invalid ID.
  - `500 Internal Server Error`: Could not delete user.

### POST /api/admin/projects/:id/members

- **Authentication:** JWT Token and Admin role required.
- **Description:** Adds a member to a project.
- **URL Parameters:**
  - `id` (integer): The ID of the project.
- **Request Body:**
  ```json
  {
    "userId": 123,
    "role": "team_developer"
  }
  ```
- **Responses:**
  - `201 Created`: Returns the new project member details (`models.ProjectMember`).
  - `400 Bad Request`: Invalid project ID or request body.
  - `409 Conflict`: User is already a member of the project.
  - `500 Internal Server Error`: Could not add member.

## WebSockets y Actualizaciones en Tiempo Real

El API provee un endpoint de WebSocket para recibir actualizaciones en tiempo real, ideal para mantener sincronizados los tableros Kanban y otras interfaces de usuario.

### Conexión

- **Endpoint:** `GET /ws`
- **Protocolo:** `wss://` (en producción) o `ws://` (en desarrollo)
- **Autenticación:** Debes proveer un token JWT válido como parámetro en la URL.
  ```
  ws://localhost:8080/ws?token=TU_TOKEN_JWT_AQUI
  ```

### Eventos del Servidor

Una vez conectado, el cliente recibirá mensajes en formato JSON con la siguiente estructura:

```json
{
  "type": "nombre_del_evento",
  "payload": { ... }
}
```

Los siguientes eventos son emitidos para las tareas:

#### `task_created`

- **Descripción:** Se emite cuando se crea una nueva tarea.
- **Payload:** El objeto completo de la nueva tarea (`models.Task`).

#### `task_status_updated`

- **Descripción:** Se emite cuando el estado (columna del Kanban) de una tarea cambia.
- **Payload:**
  ```json
  {
    "taskId": 456,
    "oldStatus": "todo",
    "newStatus": "in_progress",
    "updatedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:30:00Z"
  }
  ```

#### `task_assigned`

- **Descripción:** Se emite cuando una tarea es asignada a un usuario.
- **Payload:**
  ```json
  {
    "taskId": 456,
    "assignedTo": {
      "id": 789,
      "name": "Jane Smith"
    },
    "assignedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:31:00Z"
  }
  ```

#### `task_deleted`

- **Descripción:** Se emite cuando una tarea es eliminada.
- **Payload:**
  ```json
  {
    "taskId": 789,
    "deletedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:34:00Z"
  }
  ```